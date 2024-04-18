<template>
<div class="modal" tabindex="-1" role="dialog" ref="modal">
	<div class="modal-dialog modal-lg" role="document">
		<div class="modal-content">
			<div class="modal-body">
				<form v-on:submit.prevent="execute">
					<div class="row mb-2">
						<label class="col-sm-4 col-form-label">Mode</label>
						<div class="col-sm-8">
							<div class="form-check">
								<input class="form-ch