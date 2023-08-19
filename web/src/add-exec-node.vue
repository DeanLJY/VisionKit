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
								<div v-for="category in categories" class="tab-pa