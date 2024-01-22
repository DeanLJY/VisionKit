<template>
<div class="small-container">
	<h3>Basics</h3>
	<div class="form-group row">
		<label class="col-sm-4 col-form-label">Learning Rate</label>
		<div class="col-sm-8">
			<input v-model.number="p.LearningRate" type="text" clas