<template>
<div class="flex-container">
	<div v-if="plan && plan.length > 0">
		<h5>Execution Plan</h5>
		<table class="table table-sm w-auto">
			<thead>
				<tr>
					<th>Name</th>
					<th>Op</th>
					<th>Status</th>
				</tr>
			</thead>
			<tbody>
				<tr v-for="(vnode, idx) in plan">
					<td>{{ vnode.Name }}</td>
					<td>{{ vnode.Op }}</td>
					<td>
						<template v-if="idx < planIndex">
							Done
						</template>
						<template v-else-if="multiJob && multiJob.Done">
							<template v-if="multiJob.Error">
								Error
							</template>
							<template v-else>
								Done
							</template>
						</template>
						<template v-else-if="idx == planIndex">
							Running
						</template>
						<template v-else>
							Pending
						</template>
					</td>
				</tr>
			</tbody>
		</table>
	</div>
	<div v-if="curJob" class="flex-content">
		<component v-bind:is="'job-'+curJob.Op" v-bind:jobID="curJob.ID"></component>
	</div>
	<job-footer v-if="multiJob && multiJob.Done && !curJob" :job="multiJob"></job-footer>
</div>
</template>

<script>
import utils from './utils.js';
import JobConsoleProgress from './job-consoleprogress.vue';
import