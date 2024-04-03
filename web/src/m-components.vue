<template>
<div>
	<button type="button" class="btn btn-primary my-1" v-on:click="showAddModal">Add Component</button>
	<div class="modal" tabindex="-1" role="dialog" ref="addModal">
		<div class="modal-