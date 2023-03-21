import sys
sys.path.append('./python')
import skyhook.common as lib

import json
import numpy
import os, os.path
import requests
import skimage.io, skimage.transform

import torch
import torch.optim
import torch.utils

import skyhook.pytorch.model as model
import skyhook.pytorch.util as util

import skyhook.pytorch.dataset as skyhook_dataset
import skyhook.pytorch.augment as skyhook_augment

url = sys.argv[1]
local_port = int(sys.argv[2])
batch_size = int(sys.argv[3])

local_url = 'http://127.0.0.1:{}'.format(local_port)

# Get parameters.
resp = requests.get(local_url + '/config')
config = resp.json()

params = config['Params']
arch = config['Arch']
comps = config['Components']
datasets = config['Inputs']
parent_models = config['ParentModels']
out_dataset_id = config['Output']['ID']
train_split = config['TrainSplit']
valid_split = config['ValidSplit']

arch = arch['Params']

# overwrite parameters in arch['Components'][idx]['Params'] with parameters
# from params['Components'][idx]
if params.get('Components', None):
	overwrite_comp_params = {int(k): v for k, v in params['Components'].items()}
	for comp_idx, comp_spec in enumerate(arch['Components']):
		comp_params = {}
		if comp_spec['Params']:
			comp_params = json.loads(comp_spec['Params'])
		if overwrite_comp_params.get(comp_idx, None):
			comp_params.update(json.loads(overwrite_comp_params[comp_idx]))
		comp_spec['Params'] = json.dump