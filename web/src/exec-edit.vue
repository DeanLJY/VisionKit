<template>
<div class="el-high">
	<template v-if="node">
		<template v-if="components.includes(node.Op)">
			<component v-if="node" v-bind:is="node.Op" v-bind:node="node"></component>
		</template>
		<template v-else>
			<p>This node does not have parameters.</p>
		</template>
	</template>
</div>
</template>

<script>
import utils from './utils.js';
import CropResize from './exec-edit/cropresize.vue';
import DetectionFilter from './exec-edit/detection_filter.vue';
import ExtractPolygons from './exec-edit/extract_polygons.vue';
import GeoImageToImage from './exec-edit/geoimage_to_image.vue';
import MakeGeoImage from './exec-edit/make_geoimage.vue';
import Python from './exec-edit/python.vue';
import PytorchTrain from './exec-edit/pytorch_train.js';
import PytorchInfer from './exec-edit/pytorch_infer.vue';
import PytorchResnetTrain from './exec-edit/pytorch_resnet_train.js';
import PytorchResnetInfer from './exec-edit/pytorch_resnet_infer.vue';
import PytorchSsdTrain from './exec-edit/pytorch_ssd_train.js';
import PytorchSsdInfer from './exec-edit/pytorch_ssd_infer.vue';
import PytorchUnetTrain from './exec-edit/pytorch_unet_train.js';
import PytorchUnetInfer from './exec-edit/pytorch_unet_infer.vue';
import PytorchYolov3Train from './exec-edit/pytorch_yolov3_train.js';
import PytorchYolov3Infer from './exec-edit/pytorch_yolov3_infer.vue';
import PytorchYolov5Train from './exec-edit/pytorch_yolov5_train.js';
import PytorchYolov5Infer from './exec-edit/pytorch_yolov5_infer.vue';
import Sample from './exec-edit/sample.vue';
import SegmentationMask from './exec-edit/segmentation_mask.vue';
import SimpleTracker from './exec-edit/simple_tracker.vue';
import Split from './exec-edit/split.vue';
import ReidTracker from './exec-edit/reid_tracker.vue';
import Resample from './exec-edit/resample.vue';
import Yolov3Train from './exec-edit/yolov3_train.vue';
import Yolov3In