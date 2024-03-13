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
										<input class="form-control" type="text" v-model="path" required />
										<small class="form-text text-muted">
											<template v-if="mode == 'add'">
												The path to a file or directory from which to import files.
												The path must exist on the local disk where SkyhookML is running.
											</template>
											<template v-if="mode == 'new'">
												The path to a SkyhookML-formatted archive (.zip or directory containing db.sqlite3 and files) on the local disk where SkyhookML is running.
											</template>
										</small>
									</div>
								</div>
								<div class="row mb-2">
									<div class="col-sm-2">Symlink</div>
									<div class="col-sm-10">
										<div class="form-check">
											<input class="form-check-input" type="checkbox" v-model="symlink">
											<label class="form-check-label">
												Symlink instead of copying when possible.
											</label>
										</div>
									</div>
								</div>
								<div class="row">
									<div class="col-sm-10">
										<button type="submit" class="btn btn-primary">Import</button>
									</div>
								</div>
							</form>
						</div>
						<div class="tab-pane" id="import-upload-tab">
							<template v-if="percent === null">
								<form v-on:submit.prevent="submitUpload">
									<div class="row mb-2">
										<label class="col-sm-2 col-form-labe