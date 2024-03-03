<template>
<div class="small-container m-2">
	<template v-if="node != null">
		<p>
			Use the table below to define output splits.
			Each split is configured by the percentage of the input dataset that should end up in that split.
			The percentages across splits should add up to at most 100.
		</p>
		<table class="table">
			<thead>
				<tr>
					<th>Name</th>
					<th>Percentage</th>
					<th></th>
				</tr>
			</thead>
			<tbody>
				<tr v-for="(split, i) in params.Splits">
					<td>{{ split.Name }}</td>
					<td>{{ split.Percentage }}</td>
					<td>
						<button type="button" class="btn btn-danger" v-on:cl