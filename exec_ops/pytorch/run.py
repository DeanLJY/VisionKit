import sys
sys.path.append('./python')
import skyhook.common as lib
import skyhook.io
from skyhook.op.op import Operator

import io
import json
import torch.multiprocessing as multiprocessing
import numpy
import os, os.path
import random
import requests
import skimage.io, skimage.transform
import threading
import time
import torch

import skyhook.pytorch.model as model
import skyhook.pytorch.util as util

in_dataset_id = int(sys.argv[1])
params_arg = sys.argv[2]
# TODO: make Python /synchronized-reader endpoint accept batch size
# TODO: make batch size configurable and have auto-reduce-batch-size option
#batch_size = 16

params = json.loads(params_arg)

# For inter-task parallelism, we have:
# - Multiple ingest workers that each read data for one task at a time.
# - A single inference thread that applies the model on prepared inputs.
# - Multiple egress workers that write the outputs to local Go HTTP server.
# Ingest/egress workers operate in pairs, with same worker ID.

def ingress_worker(worker_id, params, operator, task_queue, infer_queue):
	input_options = {}
	for spec in params['InputOptions']:
		input_options[spec['Idx']] = json.loads(spec['Value'])

	while True:
		job = task_queue.get()
		request_id = job['RequestID']
		task = job['Task']

		items = [item_list[0] for item_list in task['Items']['inputs']]
		in_metadatas = []
		for item in items:
			if item['Metadata']:
				in_metadatas.append(json.loads(item['Metadata']))
			else:
				in_metadatas.append({})

		# We optimize inference over video data by handling input options in ffmpeg.
		# Here, we loop over items and update metadata to match the desired resolution.
		# We also get the framerate of the first video input (if any).
		output_defaults = {}
		for i, item in enumerate(items):
			if item['Dataset']['DataType'] != 'video':
				continue
			opt = input_options.get(i, {})
			orig_dims = in_metadatas[i]['Dims']
			new_dims = util.get_resize_dims(orig_dims, opt)
			in_metadatas[i]['Dims'] = new_dims
			item['Metadata'] = json.dumps(in_metadatas[i])
			if 'framerate' not in output_defaults:
				output_defaults['framerate'] = in_metadatas[i]['Framerate']

		rd_resp = requests.post(operator.local_url + '/synchronized-reader', json=items, stream=True)
		rd_resp.raise_for_status()

		in_dtypes = [ds['DataType'] for ds in operator.inputs]

		# Whether we have sent initialization info to the egress worker about this task yet.
		initialized_egress = False

		while True:
			# Read a batch.
			try:
				datas = skyhook.io.read_datas(rd_resp.raw, in_dtypes, in_metadatas)
			except EOFError:
				break

			# Convert datas to our input form.
			# Also get default canvas_dims based on dimensions of the first input image/video/array.
			data_len = lib.data_len(in_dtypes[0], datas[0])
			pytorch_datas = []
			for ds_idx, data in enumerate(datas):
				t = in_dtypes[ds_idx]

				if t == 'video':
					# We already handled input options by mutating the item metadata.
					# So, here, we just need to transpose.
					pytorch_data = torch.from_numpy(data).permute(0, 3, 1, 2)
				else:
					opt = input_options.get(ds_idx, {})
					cur_pytorch_datas = []
					for i in range(data_len):
						element = lib.data_index(t, data, i)
						pytorch_data = util.prepare_input(t, element, in_metadatas[ds_idx], opt)
						cur_pytorch_datas.append(pytorch_data)
					pytorch_data = util.collate(t, cur_pytorch_datas)

				if 'canvas_dims' not in output_defaults and (t == 'image' or t == 'video' or t == 'array'):
					output_defaults['canvas_dims'] = [pytorch_data.shape[3], pytorch_data.shape[2]]

				pytorch_datas.append(pytorch_data)

			# Initialize the egress worker if not done already.
			# We send this through the inference thread to synchronize with any previous closed inference jobs and such.
			if not initialized_egress:
				infer_queue.put({
					'Type': 'init',
					'WorkerID': worker_id,
					'RequestID': request_id,
					'Task': task,
					'OutputDefaults': output_defaults,
				})
				initialized_egress = True

			# Pass on to inference thread.
			infer_queue.put({
				'Type': 'infer',
				'WorkerID': worker_id,
				'Datas': pytorch_datas,
			})

		# Close the egress worker.
		# We send this through the inference thread so that any pending inference jobs for this task finish first.
		infer_queue.put({
			'Type': 'close',
			'WorkerID': worker_id,
		})

def infer_thread(in_dataset_id, params, infer_queue, egress_queues):
	device = torch.device('cuda:0')
	cpu_device = torch.device('cpu')
	model_path = 'data/items/{}/model.pt'.format(in_dataset_id)
	save_dict = torch.load(model_path)

	# overwrite parameters in save_dict['arch'] with parameters from
	# params['Components'][comp_idx]
	arch = save_dict['arch']
	if params.get('Components', None):
		overwrite_comp_params = {int(k): v for k, v in params['Components'].items()}
		for comp_idx, comp_spec in enumerate(arch['Components']):
			comp_params = {}
			if comp_spec['Params']:
				comp_params = json.loads(comp_spec['Params']