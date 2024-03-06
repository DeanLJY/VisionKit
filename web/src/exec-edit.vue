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
import ExtractPolygons from './exec-edit/extract_pol