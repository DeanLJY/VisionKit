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
				<select v-model="selected" 