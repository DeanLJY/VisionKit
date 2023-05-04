import numpy
import random

def clip(x, lo, hi):
	if x < lo:
		return lo
	elif x > hi:
		return hi
	else:
		return x

class Crop(object):
	def __init__(self, params, data_types):
		# width and height are fractions between 0 and 1
		# they indicate the relative width/height of the crop to the original dimensions
		def parse_fraction(s):
			if '/' in s:
				parts = s.split('/')
				return (int(parts[0]), int(parts[1]))
			else:
				return float(s), 1
		self.Width = parse_fraction(params['Width'])
		self.Height = parse_fraction(params['Height'])
		self.data_types = data_types

		# this augmentation should be applied in the Dataset
		self.pre_torch = True

	def forward(self, batch):
		# pick offsets to crop
		# we use same offsets for entire batch
		xoff = random.random() * (1 - self.Width[0]/self.Width[1])
		yoff = random.random() * (1 - self.Height[0]/self.Height[1])

		# transform coordinates in image to coordinates in cropped image
		def coord_transform(x, y):
			x -= xoff
			y -= yoff
			x *= self.Width[1]/self.Width[0]
			y *= self.Height[1]/self.Height[0]
			return x, y

		for i, inputs in enumerate(batch):
			if self.data_types[i] in ('image', 'video', 'array'):
				width, height = inputs[0].shape[2], inputs[0].shape[1]
				target_width = int(width * self.Width[0]) // self.Width[1]
		