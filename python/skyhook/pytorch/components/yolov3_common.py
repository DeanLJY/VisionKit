import hashlib
import os.path
import sys
import torch

class ImportContext(object):
	def __init__(self):
		self.expected_path = os.path.join('.', 'data', 'models', hashlib.sha256(b'https://github.com/ultralytics/yolov3.git').hexdigest())

	def __enter__(self):
		# from github.com/ultralytics/yolov3
		sys.path.insert(1, self.expected_path)
		return self

	def __exit__(self, exc_type, exc_value, traceback):
		# reset sys.modules
		for module_name in list(sys.modules.keys()):
			if not