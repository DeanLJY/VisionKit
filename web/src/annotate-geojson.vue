<template>
<div class="flex-container el-high">
	<template v-if="annoset != null">
		<div class="mb-2">
			<form class="row g-1 align-items-center my-1" v-on:submit.prevent="saveParams">
				<div class="col-auto">
					<label class="mx-1">Tile URL</label>
				</div>
				<div class="col-auto">
					<input type="text" class="