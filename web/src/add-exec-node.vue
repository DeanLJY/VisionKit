<template>
<div class="modal" tabindex="-1" role="dialog" ref="modal">
	<div class="modal-dialog modal-xl" role="document">
		<div class="modal-content">
			<div class="modal-body">
				<form v-on:submit.prevent="createNode">
					<div class="row mb-2">
						<label class="col-sm-2 col-form-label">Name</label>
						<div class="col-sm-10">
							<input v-model="name" class="form-control" type="text" required />
						</div>
					</div>
					<div class="row mb-2">
						<label class="col-sm-2 col-form-label">Op</label>
						<div class="col-sm-10">
							<ul class="nav nav-tabs">
								<li v-for="category in categories" class="nav-item">
									<button
										class="nav-link"
										data-bs-toggle="tab"
										:data-bs-target="'#add-node-cat-' + category.ID"
										role="tab"
										>
										{{ category.Name }}
									</button>
								</li>
							</ul>
							<div class="tab-content">
								<div v-for="category in categories" class="tab-pane" :id="'add-node-cat-' + category.ID">
									<table class="table table-row-select">
										<thead>
											<tr>
												<th>Name</th>
												<th>Description</th>
											</tr>
										</thead>
										<tbody>
											<tr
												v-for="x in category.Ops"
												:class="{selected: op != null && op.ID == x.ID}"
												v-on:click="selectOp(x)"
												>
												<td>{{ x.Name }}</td>
												<td>{{ x.Description }}</td>
											</tr>
										</tbody>
									</table>
								</div>
							</div>
						</div>
					</div>
					<template v-if="op">
						<div class="row mb-2">
							<label class="col-sm-2 col-form-label">Inputs</label>
							<div class="col-sm-10">
								<table class="table">
									<thead>
										<tr>
											<th>Name</th>
											<th>Type(s)</th>
										</tr>
									</thead>
									<tbody>
										<tr v-for="input in op.Inputs">
											<td>{{ input.Name }}</td>
											<td>
												<span v-if="input.DataTypes && input.DataTypes.length > 0">
													{{ input.DataTypes }}
												</span>
												<span v-else>
													Any
												</span>
											</td>
										</tr>
									</tbody>
								</table>
							</div>
						</div>
						<