import skyhook.common as lib
import torch
import torchvision
import torchvision.models.resnet as resnet

# Provide flag to pass pretrained=True to resnet model function.
# This makes it easy for us to develop an external script that prepares
# pre-trained parameters in a SkyhookML file dataset.
Pretrain = False

# hyperparameter constants
Means = [0.485, 0.456, 0.406]
Std = [0.229, 0.224, 0.225]

class Resnet(torch.nn.Module):
	def __init__(self, info):
		super(Resnet, self).__init__()
		example_inputs = info['example_inputs']

		# detect number of classes from int metadata
		n