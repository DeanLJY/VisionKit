<template>
<div class="flex-container el-high">
	<template v-if="annoset != null">
		<div class="mb-2">
			<form class="row g-1 align-items-center my-1" v-on:submit.prevent="saveParams">
				<div class="col-auto">
					<label class="mx-1">Tile URL</label>
				</div>
				<div class="col-auto">
					<input type="text" class="form-control mx-1" v-model="params.TileURL" />
				</div>
				<div class="col-auto">
					<button type="submit" class="btn btn-primary mx-1">Save Settings</button>
				</div>
			</form>
		</div>
		<div ref="map" style="width: 100%; height: 100%"></div>
		<div class="mt-2">
			<button type="button" class="btn btn-primary mx-1" v-on:click="saveFeatures">Save GeoJSON</button>
		</div>
	</template>
</div>
</template>

<script>
import utils from './utils.js';
import LatLonProvider from './geojson/leaflet-geosearch-provider.js';

export default {
	data: function() {
		return {
			// The annotation dataset, which must be GeoJSON type.
			annoset: null,

			// Tool parameters.
			params: {
				// The Leaflet.js URL template from which images should be fetched.
				// See https://leafletjs.com/reference-1.7.1.html#tilelayer
				TileURL: '',
			},

			// Current GeoJSON FeatureCollection.
			// On initialization: this comes from the "geojson" key in annoset.
			// Afterwards (on save): features extracted from the Leaflet instance.
			featureCollection: {
				type: "FeatureCollection",
				features: [],
			},

			// Constant item key.
			// For GeoJSON data type, we put all features into one item with this key.
			itemKey: 'geojson',

			// Currently initialized Leaflet map, if any.
			map: null,
		};
	},
	created: function() {
		const setID = this.$route.params.setid;
		utils.request(this, 'GET', '/annotate-datasets/'+setID, null, (annoset) => {
			this.annoset = annoset;

			let params;
			try {
				params = JSON.parse(this.annoset.Params);
			} catch(e) {}
			if(!params) params = {};
			if(!params.TileURL) params.TileURL = '';
			this.params = params;

			this.fetch();

			this.$store.commit('setRouteData', {
				annoset: this.annoset,
			});
		});
	},
	methods: {
		fetch: function() {
			let params = {
				format: 'json',
				t: new Date().getTime(),
			};
			utils.request(this, 'GET', '/datasets/'+this.annoset.Dataset.ID+'/items/'+this.itemKey+'/get', params, (data) => {
				if(data && data.type == 'FeatureCollection') {
					this.featureCollection = data;
				}
				this.initLeaflet();
			}, null, {error: (req, status, errorMsg) => {
				// We ignore the error if it's just that the item doesn't exist.
				// (The item not existing is expected when we first create the annotation set.)
				if(req && req.responseText && req.responseText.includes('no such item')) {
					this.initLeaflet();
					return;
				}
				this.$globals.app.setError(errorMsg);
			}});
		},
		initLeaflet: function() {
			// If TileURL is not configured yet, there's not much point to display anything.
			if(!this.params.TileURL) {
				return;
			}

			// Create GeoJSON layer in Leafle