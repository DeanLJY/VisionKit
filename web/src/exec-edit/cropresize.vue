<template>
<div class="small-container m-2">
	<template v-if="node != null">
		<div class="form-group row">
			<label class="col-sm-2 col-form-label">X Offset</label>
			<div class="col-sm-10">
				<input v-model.number="start[0]" type="text" class="form-control">
			</div>
		</div>
		<div class="form-group row">
			<label class="col-sm-2 col-form-label">Y Offset</label>
			<div class="col-sm-10">
				<input v-model.number="start[1]" type="text" class="form-control">
			</div>
		</div>
		<div class="form-group row">
			<label class="col-sm-2 col-form-label">Crop Width</label>
			<div class="col-sm-10">
				<input v-model.number="cropDims[0]" type="text" class="form-control">
			</div>
		</div>
		<div class="form-group row">
			<label class="col-sm-2 col-form-label">Crop Height</label>
			<d