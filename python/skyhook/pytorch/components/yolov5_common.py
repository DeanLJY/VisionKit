import hashlib
import os.path
import sys
import torch

class ImportContext(object):
	def __init__(self):
		self.expected_path = os.path.join('.', 'data', 'models'