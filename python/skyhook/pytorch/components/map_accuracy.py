import numpy
import torch

# Mean average precision metric for object detection.
# Detection format: (cls, xyxy, conf)

# Group detections by image, and move to numpy.
def group_by_image(raw_detections):
	counts = raw_detections['count