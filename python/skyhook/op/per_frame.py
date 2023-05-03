from skyhook.op.op import Operator

import skyhook.common as lib
import skyhook.io

import io
import json
import requests

class PerFrameOperator(Operator):
	def __init__(self, meta_packet):
		super(PerFrameOperator, self).__init__(meta_packet)
		# Function must be set after initialization.
		self.f = None

	def apply(self, task):
		# Use SynchronizedReader to read the input items chunk-by-chunk.
		items = [item_list[0] for item_list in task['Items']['inputs']]
		rd_resp = requests.post(self.local_url + '/synchronized-reader', json=items, stream=True)
		rd_resp.raise_for_status()

		in_dtypes = [ds['DataType'] for ds in self.inputs]
		out_dtypes = [ds['DataType'] for ds in self.outputs]

		in_metadatas = []
		for ds_idx, item in enumerate(items):
			in_metadatas.append(lib.decode_metadata(self.inputs[ds_idx], item))

		# Define generator function that will run self.f on each element of sequence data.
		def gen():
			sent_meta = False

			while True:
				# Read a chunk.
				try:
					datas = skyhook.io.read_datas(rd_resp.raw, in_dtypes, in_metadatas)
				except EOFError:
					break

				# Collect output chunk by running self.f on each element.
				input_len = lib.data_len(in_dtypes[0], datas[0])
				out_datas = []
				out_metadatas = []
				for i in range(input_len):
					cur_inputs = [lib.data_index(in_dtypes[ds_idx], data, i) for ds_idx, data in enumerate(datas)]
					cur_inputs = [{'Data': data, 'Metadata': in_metadatas[ds_idx]} for ds_idx, data in enumerate(cur_inputs)]
					cur_outputs = self.f(*cur_inputs)
					if not isinstance(cur_outputs, tuple):
						cur_outputs = (cur_outputs,)

					cur_datas = []
					cur_metadatas = []
					for arg in cur_outputs:
						if isinstance(arg, dict) and 'Data' in arg:
							cur_datas.append(arg['Data'])
							cur_metadatas.append(arg['Metadata'])
						else:
							cur_datas.append(arg)
							cur_metadatas.append({})
					out_datas.append(cur_datas)
					out_metadatas.appen