<template>
<div class="small-container m-2">
	<template v-if="node != null">
		<div class="form-group row">
			<label class="col-sm-2 col-form-label">Categories</label>
			<div class="col-sm-10">
				<table class="table">
					<tbody>
						<tr v-for="(category, i) in categories">
							<td>{{ category }}</td>
							<td>
								<button type="button" class="btn btn-danger" v-on:click="removeCategory(i)">Remove</button>
							</td>
						</tr>
						<tr>
							<td>
								<input v-model="addCategoryInput" type="text" class="form-control">
							</td>
							<td>
								<button type="button" class="btn btn-primary" v-on:click="addCategory">Add</button>
							</td>
						</tr>
					</tbody>
				</table>
			</div>
		</div>
		<div class="form-group row">
			<label class="col-sm-2 col-form-label">Score Threshold</label>
			<div class="col-sm-10">
				<input v-model.number="score" type="text" class="form-control">
			</div>
		</div>
		<button v-on:click="save" type="button" class="btn btn-primary">Save</button>
	</template>
</div>
</template>

<script>
impor