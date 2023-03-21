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
		comp_spec['Params'] = json.dumps(comp_params)

device = torch.device('cuda:0')
#device = torch.device('cpu')

# get train and val Datasets
print('loading datasets')
dataset_provider = skyhook_dataset.providers[params['Dataset']['Op']]
dataset_params = json.loads(params['Dataset']['Params'])
train_set, val_set = dataset_provider(url, datasets, dataset_params, train_split, valid_split)
datatypes = train_set.get_datatypes()

# get data augmentation steps
# this is a list of objects that provide forward() function
# we will apply the forward function on batches from DataLoader
print('loading data augmentations')
ds_augments = []
torch_augments = []
for spec in params['Augment']:
	cls_func = skyhook_augment.augmentations[spec['Op']]
	obj = cls_func(json.loads(spec['Params']), datatypes)
	if obj.pre_torch:
		ds_augments.append(obj)
	else:
		torch_augments.append(obj)

train_set.set_augments(ds_augments)
val_set.set_augments(ds_augments)

# apply data augmentation on validation set
# this is because some augmentations are random but we want a consistent validation set
# here we assume the validation set fits in system memory, but not necessarily GPU memory
# so we apply augmentation on CPU, whereas during training we will apply on GPU

train_params = json.loads(params['Train']['Params'])

print('preparing validation set')
val_loader = torch.utils.data.DataLoader(
	val_set,
	batch_size=batch_size,
	num_workers=4,
	collate_fn=val_set.collate_fn,
	# drop last unless we'd end up with 0 batches
	drop_last=len(val_set) > batch_size
)
val_batches = []
for batch in val_loader:
	for obj in torch_augments:
		batch = obj.forward(batch)
	val_batches.append(batch)

'''
batch = val_batches[0]
for i in range(32):
	im = batch[0][i, :, :, :].cpu().numpy().transpose(1, 2, 0)
	prefix = sum(batch[1]['counts'][0:i])
	detections = batch[1]['detections'][prefix:prefix+batch[1]['counts'][i]]
	for d in detections:
		cls, sx, sy, ex, ey = d
		sx = int(sx*im.shape[1])
		sy = int(sy*im.shape[0])
		ex = int(ex*im.shape[1])
		ey = int(ey*im.shape[0])
		im[sy:sy+2, sx:ex, :] = [255, 0, 0]
		im[ey-2:ey, sx:ex, :] = [255, 0, 0]
		im[sy:ey, sx:sx+2, :] = [255, 0, 0]
		im[sy:ey, ex-2:ex, :] = [255, 0, 0]
	skimage.io.imsave('/home/ubuntu/vis/{}.jpg'.format(i), im)
'''

'''
batch = val_batches[0]
for i in range(32):
	im1 = batch[0][i, :, :, :].cpu().numpy().transpose(1, 2, 0)
	im2 = (batch[1][i, 0, :, :].cpu().numpy() > 0).astype('uint8')*255
	skimage.io.imsave('/home/ubuntu/vis/{}_im.jpg'.format(i), im1)
	skimage.io.imsave('/home/ubuntu/vis/{}_mask.png'.format(i), im2)
'''

print('initialize model')
train_loader = torch.utils.data.DataLoader(
	train_set,
	batch_