import io
import json
import math
import numpy
import os
import os.path
import skimage.io
import struct
import sys

import skyhook.ffmpeg as ffmpeg
import skyhook.io

def eprint(s):
	sys.stderr.write(str(s) + "\n")
	sys.stderr.flush()

# sometimes JSON that we input ends up containing null (=> None) entries instead of list
# this helper restores lists where lists are expected
def non_null_list(l):
	if l is None:
		return []
	return l

def data_index(t, data, i):
	return data[i]

# stack a bunch of individual data (like data_index output)
def data_stack(t, datas):
	if t == 'image' or t == 'video' or t == 'array':
		return numpy.stack(datas)
	else:
		return datas

# stack a bunch of regular data
# this fails for non-sequence data, unless len(datas)==1, in which case it simply returns the data
def data_concat(t, datas):
	if len(datas) == 1:
		return datas[0]

	if t == 'image' or t == 'video' or t == 'array':
		return numpy.concatenate(datas, axis=0)
	else:
		return [x for data in datas for x in data]

def data_len(t, data):
	return len(data)

def decode_metadata(dataset, item):
	metadata = {}
	if dataset['Metadata']:
		metadata.u