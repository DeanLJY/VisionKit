<template>
<div class="small-container m-2">
	<template v-if="node != null">
		<div class="form-group row">
			<label class="col-sm-2 col-form-label">Architecture</label>
			<div class="col-sm-10">
				<select v-model="params.archID" class="form-select">
					<template v-for="arch in archs">
						<option :key="arch.ID" :value="arch.ID">{{ arch.ID }}</option>
					</template>
				</select>
			</div>
		</div>
		<template v-if="arch">
			<div class="form-group row" v-if="parents.length > 0">
				<label class="col-sm-2 col-form-label">Input Options</label>
				<div class="col-sm-10">
					<table class="table">
						<tbody>
							<tr v-for="(spec, i) in params.inputOptions">
								<td>{{ parents[spec.Idx].Name }}</td>
								<td>{{ spec.Value }}</td>
								<td>
									<button type="button" class="btn btn-danger" v-on:click="removeInput(i)">Remove</button>
								</td>
							</tr>
							<tr>
								<td>
									<select v-model="addForms.inputIdx" class="form-select">
										<template v-for="(parent, parentIdx) in parents">
											<option :value="parentIdx">{{ parent.Name }} ({{ parent.DataType }})</option>
										</template>
									</select>
								</td>
								<td>
									<input class="form-control" type="text" v-model="addForms.inputOptions" />
								</td>
								<td>
	