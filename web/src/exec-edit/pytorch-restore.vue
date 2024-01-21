<template>
<div class="small-container">
	<template v-for="(name, idx) in parentNames">
		<h4>{{ name }}</h4>
		<div class="form-group row">
			<label class="col-sm-4 col-form-label">Source Prefix</label>
			<div class="col-sm-8">
				<input v-model.number="restore[idx].SrcPrefix" type="text" class="form-control" @change="update">
				<small class="form-text text-muted">
					Prefix of data in the parent model to load.
				</small>
			</div>
		</div>
		<div class="form-group row">
			<label class="col-sm-4 col-form-label">Destination Prefix</label>
			<div class="col-sm-8">
				<input v-model.number="restore[idx].DstPrefix" type="text" class="form-control" @change="update">
				<small class="form-text text-muted">
					Prefix of data in this model to restore the saved parameters to.
				</small>
			</div>
		</div>
		<div class="form-group row">
			<label class="col-sm-4 col-form-label">Skip Prefixes</label>
			<div class="col-sm-8">
				<input v-model.number="restore[idx].SkipPrefixes" type="text" class="form-control" @change="update">
				<small class="form-text text-muted">
					Comma-separa