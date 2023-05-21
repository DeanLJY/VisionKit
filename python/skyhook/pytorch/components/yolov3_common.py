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
			if not hasattr(sys.modules[module_name], '__file__'):
				continue
			fname = sys.modules[module_name].__file__
			if fname is None:
				continue
			if not fname.startswith(self.expected_path):
				continue
			del sys.modules[module_name]
		sys.path.remove(self.expected_path)

# process skyhook detections into yolov3 target detections
def process_targets(targets):
	# first extract detection counts per image in the batch, and the boxes
	if 'points' in targets:
		# shape type
		counts = targets['counts']
		cls_labels = targets['infos'][:, 0].float()
		boxes = targets['points'].reshape(-1, 4)
		# need to make sure that first point is smaller than second point
		boxes = torch.stack([
			torch.minimum(boxes[:, 0], boxes[:, 2]),
			torch.minimum(boxes[:, 1], boxes[:, 3]),
			torch.maximum(boxes[:, 0], boxes[:, 2]),
			torch.maximum(boxes[:, 1], boxes[:, 3]),
		], dim=1)
	elif 'detections' in targets:
		# detection type
		counts = targets['counts']
		cls_labels = targets['detections'][:, 0]
		boxes = targets['detections'][:, 1:5]

	# xyxy -> xywh
	boxes = torch.stack([
		(boxes[:, 0] + boxes[:, 2]) / 2,
		(boxes[:, 1] + boxes[:, 3]) / 2,
		boxes[:, 2] - boxes[:, 0],
		boxes[:, 3] - boxes[:, 1],
	], dim=1)

	# output: list of detections with:
	# - first column indicating image index
	# - second column indicating class index
	ind