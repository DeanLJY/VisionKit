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
		rd_resp.raise_