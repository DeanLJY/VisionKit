import sys
sys.path.append('./python')

import json
import numpy
import os, os.path
import random
import requests
import skimage.io, skimage.transform
import struct
import sys

import torch

import skyhook.pytorch.model as model
import skyhook.pytorch.util as util

in_dataset_id = int(sys.argv[1])

device = torch.device('cuda:0')
#device = torch.device('cpu')
model_path = 'data/items/{}/model.pt'.format(in_dataset_id)
save_dict = torch.load(model_path)
example_inputs = save_dict['example_inputs']
util.inputs_to_device(example_inputs, device)
net = model.Net(save_