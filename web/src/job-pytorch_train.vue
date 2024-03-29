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
		<job-footer :job="job"></job-footer>
	</div>
</div>
</template>

<script>
import utils from './utils.js';
import JobConsole from './job-console.vue';
import JobFooter from './job-footer.vue';

export default {
	components: {
		'job-console': JobConsole,
		'job-footer': JobFooter,
	},
	data: function() {
		return {
			job: null,
			modelState: null,
			lines: [],
			chart: null,
		};
	},
	props: ['jobID'],
	created: function() {
		this.fetch();
		this.interval = setInterval(this.fetch, 1000);
	},
	destroyed: function() {
		clearInterval(this.interval);
	},
	methods: {
		fetch: function() {
			utils.request(this, 'POST', '/jobs/'+this.jobID+'/state', null, (response) => {
				this.job = response.Job;
				let state;
				try {
					state = JSON.parse(response.State);
				} catch(e) {}
				if(!state) {
					return;
				}
				let metadata = null;
				try {
					metadata = JSON.parse(state.Datas.node);
				} catch(e) {}
				this.updateChart(metadata);
				this.lines = state.Lines;
			});
		},
		updateChart: function(modelState) {
			if(!modelState || !modelState.TrainLoss || modelState.TrainLoss.length == 0) {
				return;
			}
			let n = modelState.TrainLoss.length;
			let prevN = 0;
			if(this.modelState) {
				prevN = this.modelState.TrainLoss.length;
			}
			if(prevN == n) {
				return;
			}
			if(!this.chart) {
				let labels = [];
				for(let i = 