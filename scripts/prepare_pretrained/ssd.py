import hashlib
import json
import numpy
import os.path
import sys
import torch
import yaml

sys.path.append('./python/')
import skyhook.pytorch.model as model
import skyhook.pytorch.util as util

mode = sys.argv[1]
in_fname = sys.argv[2]
out_fname = sys.argv[3]

device = torch.device('cpu')
ssd_path = os.path.joi