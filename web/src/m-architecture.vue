
<template>
<div class="small-container m-2">
	<template v-if="arch != null">
		<div class="form-group row">
			<label class="col-sm-2 col-form-label"># Inputs</label>
			<div class="col-sm-10">
				<input v-model="numInputs" type="text" class="form-control">
			</div>
		</div>
		<div class="form-group row">
			<label class="col-sm-2 col-form-label"># Targets</label>
			<div class="col-sm-10">
				<input v-model="numTargets" type="text" class="form-control">
			</div>
		</div>
		<div class="form-group row">
			<label class="col-sm-2 col-form-label">Components</label>
			<div class="col-sm-10">
				<template v-for="(compSpec, compIdx) in components">
					<h3>
						Component #{{ compIdx }}: {{ compSpec.ID }}
						<button type="button" class="btn btn-sm btn-danger" v-on:click="removeComponent(compIdx)">Remove</button>
					</h3>
					<p>Inputs:</p>
					<table class="table">
						<tbody>
							<tr v-for="(inp, i) in compSpec.Inputs">
								<td>{{ inp.Type }}</td>
								<td>
									<template v-if="inp.Type == 'dataset'">Dataset {{ inp.DatasetIdx }}</template>
									<template v-else-if="inp.Type == 'layer'">Component #{{ inp.ComponentIdx }} / Layer {{ inp.Layer }}</template>
								</td>
								<td>
									<button type="button" class="btn btn-danger" v-on:click="removeInput(compIdx, i)">Remove</button>
								</td>
							</tr>
						</tbody>
					</table>
					<button type="button" class="btn btn-primary" v-on:click="showAddInputModal(compIdx, 'inputs')">Add Input</button>
					<p>Targets:</p>
					<table class="table">
						<tbody>
							<tr v-for="(inp, i) in compSpec.Targets">
								<td>{{ inp.Type }}</td>
								<td>
									<template v-if="inp.Type == 'dataset'">Dataset {{ inp.DatasetIdx }}</template>
									<template v-else-if="inp.Type == 'layer'">Component #{{ inp.ComponentIdx }} / Layer {{ inp.Layer }}</template>
								</td>
								<td>
									<button type="button" class="btn btn-danger" v-on:click="removeTarget(compIdx, i)">Remove</button>
								</td>
							</tr>
						</tbody>
					</table>
					<button type="button" class="btn btn-primary" v-on:click="showAddInputModal(compIdx, 'targets')">Add Target</button>
					<p>Parameters:</p>
					<textarea v-model="compSpec.Params" class="form-control" rows="5"></textarea>
				</template>
				<hr />
				<form class="row g-1 align-items-center" v-on:submit.prevent="addComponent">
					<div class="col-auto">
						<select class="form-select" v-model="addComponentForm.componentID">
							<option v-for="comp in comps" :key="comp.ID" :value="comp.ID">{{ comp.ID }}</option>
						</select>
					</div>
					<div class="col-auto">
						<button type="submit" class="btn btn-primary my-1 mx-1">Add Component</button>
					</div>
				</form>
			</div>
		</div>
		<template v-if="addInputModal">
			<m-architecture-input-modal v-bind:components="parentComponentList(addInputModal.componentIdx)" v-on:success="addInput($event)"></m-architecture-input-modal>
		</template>
		<template v-for="{label, list, form} in lossAndScore">
			<div class="form-group row">
				<label class="col-sm-2 col-form-label">{{ label }}</label>
				<div class="col-sm-10">
					<table class="table">
						<thead>
							<tr>
								<th>Component</th>
								<th>Layer</th>
								<th></th>
							</tr>
						</thead>
						<tbody>
							<tr v-for="(spec, i) in list">
								<td>
									Component #{{ spec.ComponentIdx }}<template v-if="getComponent(spec.ComponentIdx)">: {{ getComponent(spec.ComponentIdx).ID }}</template>
								</td>
								<td>{{ spec.Layer }}</td>
								<td>{{ spec.Weight }}</td>
								<td>
									<button type="button" class="btn btn-danger" v-on:click="removeLoss(list, i)">Remove</button>
								</td>
							</tr>
							<tr>
								<td>
									<select v-model="form.componentIdx" class="form-select">
										<template v-for="(compSpec, compIdx) in components">
											<option v-if="compSpec.ID in comps" :key="compIdx" :value="compIdx">Component #{{ compIdx }}: {{ compSpec.ID }}</option>
										</template>
									</select>
								</td>
								<td>
									<template v-if="getComponent(form.componentIdx)">
										<select v-model="form.layer" class="form-select">
											<template v-for="layer in getComponent(form.componentIdx).Params.Losses">
												<option :key="layer" :value="layer">{{ layer }}</option>
											</template>
										</select>
									</template>
								</td>
								<td>
									<input class="form-control" type="text" v-model.number="form.weight" />
								</td>
								<td>
									<button type="button" class="btn btn-primary" v-on:click="addLoss(list, form)">Add</button>
								</td>
							</tr>
						</tbody>
					</table>
				</div>
			</div>
		</template>
		<button v-on:click="save" type="button" class="btn btn-primary">Save</button>
	</template>
</div>
</template>
