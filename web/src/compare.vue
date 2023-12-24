
<template>
<div>
	<template v-if="viewingItems != null">
		<render-item v-bind:dataType="dataset1.DataType" v-bind:item="viewingItems[0]"></render-item>
		<render-item v-bind:dataType="dataset2.DataType" v-bind:item="viewingItems[1]"></render-item>
		<div class="form-row align-items-center">
			<div class="col-auto">
				<button v-on:click="viewNext(-1)" type="button" class="btn btn-primary">Prev</button>
			</div>
			<div class="col-auto">
				<span>{{ viewingItems[0].Key }}</span>
				<span>({{ curIndex }} of {{ items.length }})</span>
			</div>
			<div class="col-auto">
				<button v-on:click="viewNext(1)" type="button" class="btn btn-primary">Next</button>
			</div>
		</div>
	</template>
</div>
</template>

<script>
import utils from './utils.js';
import RenderItem from './render-item.vue';
