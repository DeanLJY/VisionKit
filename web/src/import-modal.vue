<template>
<span>
	<button type="button" class="btn btn-primary" v-on:click="click">
		<template v-if="mode == 'add'">
			Import Data
		</template>
		<template v-else-if="mode == 'new'">
			Import SkyhookML Dataset
		</template>
	</button>
	<div class="modal" tabindex="-1" role="dialog" ref="modal">
		<div class="modal-dialog" role="document">
			<div class="modal-content">
				<div class="modal-body">
					<ul class="nav nav-tabs mb-2">
						<li class="nav-item">
							<button class="nav-link active" data-bs-toggle="tab" data-bs-target="#import-local-tab" role="tab">Local</button>
						</li>
						<li class="nav-item">
							<button class="nav-link" data-bs-toggle="tab" data-bs-target="#import-upload-tab" role="tab">Upload</button>
						</li>
						<li class="nav-item">
							<button class="nav-link" data-bs-toggle="tab" data-bs-target="#import-url-tab" role="tab">URL</button>
						</li>
					</ul>
					<div class="tab-content">
						<div class="tab-pane show active" id="import-local-tab">
							<form v-on:submit.prevent="submitLocal">
								<div class="row mb-2">
									<label class="col-sm-2 col-form-label">Path</label>
									<div class="col-sm-10">
										<input class="form-control" type="text" v-model="path