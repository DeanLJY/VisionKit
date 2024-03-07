<template>
<div>
	<template v-if="this.isVariable">
		<table class="table table-sm w-auto">
			<thead>
				<tr><th colspan="2">{{ input.Name }}</th></tr>
			</thead>
			<tbody>
				<tr v-for="(parent, i) in node.Parents[input.Name]" :key="i">
					<template v-if="parent.Type == 'd'">
						<td>Dataset: {{ datasets[parent.ID].Name }}</td>
					</template>
					<template v-else-if="parent.Type == 'n'">
						<td>Node: {{ nodes[parent.ID].Name }}[{{ parent.Name }}]</td>
					</template>
					<td><button type="button" class="btn btn-danger btn-sm" v-on:click="$emit('remove', i)">Remove</button></td>
				</tr>
				<tr>
					<td>
						<select v-model="selected" class="form-select">
							<template v-for="(label, key) in options">
								<option :value="key">{{ label }}</option>
							</template>
						</select>
					</td>
					<td><button type="button" class="btn btn-success btn-sm" v-on:click="add">Add</button></td>
				</tr>
			</tbody>
		</table>
	</template>
	<template v-else>
		<div class="d-flex">
			<form class="d-flex align-items-center">
				<label class="mx-2">{{ input.Name }}</label>
				<select v-model="selected" @change="parentChanged" class="form-select form-select-sm mx-2">
					<option value="">None</option>
					<template v-for="(label, key) in options">
						<option :value="key" :key="key">{{ label }}</option>
					</template>
				</select>
			</form>
		</div>
	</template>
</div>
</template>

<script>
import utils from './utils.js';

export default {
	data: function() {
		return {
			isVariable: false,

			// Current parents.
			parents: [],

			// Options for this parent, along with the currently selected option.
			// This option is for:
			// (1) Setting or unsetting the single parent if this input is non-variable
			// (2) Adding a new parent if this input is variable
			selected: '',
			options: {},

			// Map from option names to the corresponding ExecParent object.
			optionToObj: {},
		};
	},
	props: [
		'node', 'input', 'nodes', 'datasets',
	],
	created: function() {
		this.parents = this.node.Parents[this.input.Name];
		this.isVariable = this.input.Variable || this.parents.length > 1;

		// helper function that decides whether a given data type is acceptable for this input
		let dataTypeSet = null;
		if(this.input.DataTypes && this.input.DataTypes.length > 0) {
			dataTypeSet = {};
			this.input.DataTypes.forEach((dt) => {
				dataTypeSet[dt] = true;
			});
		}
		let isTypeOK = (dt) => {
			return !dataTypeSet || dataTypeSet[dt];
		};
