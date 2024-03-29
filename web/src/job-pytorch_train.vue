<template>
<div class="flex-container">
	<div class="chartjs-container el-50h">
		<canvas ref="layer"></canvas>
	</div>
	<div class="el-50h flex-container">
		<job-console :lines="lines"></job-console>
		<div v-if="job && !job.Done" class="mb-2">
			<button class="btn btn-warning" v-on:click="stopTraining" data-bs-toggle="tooltip" title="Terminate the job, and mark the currently saved model as completed.">Stop Training and Mark Done</button>
		</div>
		<job-footer :job="job"></job-fo