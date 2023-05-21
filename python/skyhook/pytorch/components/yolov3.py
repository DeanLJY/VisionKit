import os.path
import skyhook.common as lib
import torch
import yaml

import skyhook.pytorch.components.yolov3_common as yolov3_common

def M(info):
	with yolov3_common.ImportContext() as ctx:
		import utils.general
		import utils.loss
		import models.yolo

		class Yolov3(torch.nn.Module):
			def __init__(self, info):
				super(Yolov3, self).__init__()
				self.infer = info['infer']
				detection_metadata = info['metadatas'][1]
				if detection_metadata and 'Categories' in detection_metadata:
					self.categories = detection_metadata['Categories']
				else:
					self.categories = ['object']
				self.nc = len(self.categories)

				# e.g. 'yolov3', 'yolov3-tiny', 'yolov3-spp'
				self.mode = info['params'].get('mode', 'yolov3')

				if self.infer:
					default_confidence_threshold = 0.1
				else:
					default_confidence_threshold = 0.01
				self.confidence_threshold = info['params'].get('confidence_threshold', default_confidence_threshold)
				self.iou_threshold = info['params'].get('iou_threshold', 0.5)

				lib.eprint('yolov3: set nc={}, mode={}, conf={}, iou={}'.format(self.nc, self.mode, self.confidence_threshold, self.iou_threshold))

				self.model = models.yolo.Model(cfg=os.path.join(ctx.expected_path, 'models', '{}.yaml'.format(self.mode)), nc=self.nc)
				self.model.nc = self.nc
				with open(os.path.join(ctx.expected_path, 'data', 'hyp.scratch.yaml'), 'r') as f:
					hyp = yaml.load(f, Loader=yaml.FullLoader)
				self.model.hyp =