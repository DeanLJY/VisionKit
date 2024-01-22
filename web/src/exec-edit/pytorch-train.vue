<template>
<div class="small-container">
	<h3>Basics</h3>
	<div class="form-group row">
		<label class="col-sm-4 col-form-label">Learning Rate</label>
		<div class="col-sm-8">
			<input v-model.number="p.LearningRate" type="text" class="form-control" @change="update">
		</div>
	</div>
	<div class="form-group row">
		<label class="col-sm-4 col-form-label">Optimizer</label>
		<div class="col-sm-8">
			<select v-model="p.Optimizer" class="form-select" @change="update">
				<option value="adam">Adam</option>
			</select>
		</div>
	</div>
	<div class="form-group row">
		<label class="col-sm-4 col-form-label">Batch Size</label>
		<div class="col-sm-8">
			<input v-model.number="p.BatchSize" type="text" class="form-control" @change="update">
		</div>
	</div>
	<div class="form-group row">
		<label class="col-sm-4 col-form-label">Auto Batch Size</label>
		<div class="col-sm-8">
			<div class="form-check">
				<input class="form-check-input" type="checkbox" v-model="p.AutoBatchSize" @change="update">
				<label class="form-check-label">
					Automatically reduce the batch size if we run out of GPU memory.
				</label>
			</div>
		</div>
	</div>

	<h3>Stop Condition</h3>
	<div class="form-group row">
		<label class="col-sm-4 col-form-label">Max Epochs</label>
		<div class="col-sm-8">
			<input v-model.number="p.StopCondition.MaxEpochs" type="text" class="form-control" @change="update">
			<small class="form-text text-muted">
				Stop training after this many epochs. Leave 0 to disable this stop condition.
			</small>
		</div>
	</div>
	<div class="form-group row">
		