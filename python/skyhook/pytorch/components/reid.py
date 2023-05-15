import skyhook.common as lib
import torch

class Reid(torch.nn.Module):
	def __init__(self, info):
		super(Reid, self).__init__()

		conv_layers = [
			torch.nn.Conv2d(3, 32, 4, padding=(1, 1), stride=2), # 32