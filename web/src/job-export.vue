<template>
<div class="el-high">
	<template v-if="!job || !job.Done || job.Error">
		<job-consoleprogress :jobID="jobID"></job-consoleprogress>
	</template>
	<templa