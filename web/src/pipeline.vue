<template>
<div class="flex-container">
	<div class="my-2">
		<div class="float-start">
			<h3>Pipeline</h3>
		</div>
		<div class="btn-group float-end">
			<button
				class="btn btn-outline-secondary shadow-none"
				:class="{active: mode == 'table'}"
				v-on:click="mode = 'table'"
				>
				Table View
			</button>
			<button
				c