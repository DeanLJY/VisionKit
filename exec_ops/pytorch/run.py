import sys
sys.path.append('./python')
import skyhook.common as lib
import skyhook.io
from skyhook.op.op import Operator

import io
import json
import torch.multiprocessing as multiprocessing
import numpy
import os, os.path
import random
import requests
import skimage.io, skimage.transform
import threading
import time
import torch

import skyhook.pytorch.model as model
import skyhook.pytorch.util as util

in_dataset_id = int(sys.argv[1])
params_arg = sys.argv[2]
# TODO: make Python /synchronized-reader endpoint accept batch size
# TODO: make batch size configurable and have auto-reduce-batch-size option
#batch_size = 16

params = json.loads(params_arg)

# For inter-task parallelism, we have:
# - Multiple ingest workers that each read data for one task at a time.
# - A single inference thread that applies the model on prepared inputs.
# - Multiple egress workers that write the outputs to local Go HTTP se