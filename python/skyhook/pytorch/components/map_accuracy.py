import numpy
import torch

# Mean average precision metric for object detection.
# Detection format: (cls, xyxy, conf)

# Group detections by image, and move to numpy.
def group_by_image(raw_detections):
	counts = raw_detections['counts'].detach().cpu().numpy()
	raw_dlist = raw_detections['detections'].detach().cpu().numpy()
	prefix_sum = 0
	dlists = []
	for count in counts:
		dlists.append(raw_dlist[prefix_sum:prefix_sum+count])
		prefix_sum += count
	return dlists

# Group detections by category IDs.
# But only keep categories with at least one ground truth detection.
def group_by_category(pred, gt):
	categories = {}
	for d in gt:
		category_id = int(d[0])
		if category_id not in categories:
			categories[category_id] = ([], [])
		categories[category_id][1].append(d[1:5])
	for d in pred:
		category_id = int(d[0])
		if category_id not in categories:
			continue
		categories[category_id][0].append(d[1:6])
	return categories

# Return intersection-over-union between boxes.
# Returns IOU, where IOU[i, j] is IOU between pred[i] and gt[j]
def get_iou(pred, gt):
	def box_area(box):
		return (box[:, 2] - box[:, 0]) * (box[:, 3] - box[:, 1])

	area1 = box_area(pred)
	area2 = box_area(gt)

	intersect_area = (numpy.minimum(pred[:, None, 2:4], gt[:, 2:4]) - numpy.maximum(pred[:, None, 0:2], gt[:, 0:2]))
	intersect_area = numpy.maximum(intersect_area, 0)
	intersect_area = numpy.prod(intersect_area, axis=2)
	union_area = area1[:, None] + area2 - intersect_area
	return intersect_area / union_area

def compute_ap(recall, precision):
	# Appe