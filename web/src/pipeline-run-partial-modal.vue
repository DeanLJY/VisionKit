<template>
<div class="modal" tabindex="-1" role="dialog" ref="modal">
	<div class="modal-dialog modal-lg" role="document">
		<div class="modal-content">
			<div class="modal-body">
				<form v-on:submit.prevent="execute">
					<div class="row mb-2">
						<label class="col-sm-4 col-form-label">Mode</label>
						<div class="col-sm-8">
							<div class="form-check">
								<input class="form-check-input" type="radio" v-model="mode" value="random">
								<label class="form-check-label">Random: Compute a fixed number of random outputs.</label>
							</div>
							<div class="form-check">
								<input class="form-check-input" type="radio" v-model="mode" value="dataset">
								<label class="form-check-label">Dataset: Compute only outputs with keys matching those in another dataset.</label>
							</div>
							<div class="form-check">
								<input class="form-check-input" type="radio" v-model="mode" value="direct">
								<label class="form-check-label">Direct: Compute outputs matching a specified list of keys.</label>
							</div>
						</div>
					</div>
					<template v-if="mode == 'random'">
						<div class="form-group row">
							<label class="col-sm-4 col-form-label">Count</label>
							<div class="col-sm-8">
								<input v-model.number="count" type="text" class="form-control">
								<small class="form-text text-muted">
									The number of output items to compute.
								</small>
							</div>
						</div>
					</template>
					<template v-if="mode == 'dataset'">
						<div class="form-group row">
							<label class="col-sm-4 col-form-label">Dataset</label>
							<div class="col-sm-8">
								<select v-model="optionIdx" class="form-select" required>
									<template v-for="(opt, idx) in options">
										<optio