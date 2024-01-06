<template>
<div class="small-container">
	<div class="form-group row">
		<label class="col-sm-2 col-form-label">Validation Percentage</label>
		<div class="col-sm-10">
			<input v-model.number="valPercent" type="text" class="form-control" @change="update">
			<small class="form-text text-muted">
				Use this percentage of the input data for validation. The rest will be used for training.
			</small>
		</div>
	</div>
	<h3>Input Options</h3>
	<template v-for="(parent, idx) in parents">
		<h4>{{ parent.Name }} ({{ parent.DataType }})</h4>
		<template v-if="['image', 'video', 'array'].includes(parent.DataType)">
			<select-input-size v-model="inputOptions[idx]" @change="update"></select-input-size>
		</template>
	</template>
</div>
</template>

<script>
import utils from '../utils.js';
import SelectInputSize from './select-input-size.vue';

export default {
	components: {
		'select-input-size': SelectInputSize,
	},
	data: function() {
		return {
			parents: [],
			inputOptions: [],
			valPercent: 20,
		};
	},
	props: ['node', 'value'],
	created: function() {
		try {
			let s = JSON.parse(this.value);
			if(s.InputOption