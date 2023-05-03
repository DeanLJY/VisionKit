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

		in_dtypes = [ds['DataType'] for ds in self.inpu