<template>
<div class="small-container m-2">
	<template v-if="node != null">
		<div class="form-group row">
			<label class="col-sm-2 col-form-label">Fraction</label>
			<div class="col-sm-10">
				<input v-model="fraction" type="text" class="form-control">
				<small class="form-text text-muted">
					Re-sample input data by this fraction