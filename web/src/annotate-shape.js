import utils from './utils.js';
import AnnotateGenericUI from './annotate-generic-ui.js';

export default AnnotateGenericUI({
	data: function() {
		return {
			params: null,
			shapes: null,

			// current category to use for labeling shapes
			category: '',

			// index of currently selected shape, if any
			selectedIdx: null,

			keyupHandler: null,
			resizeObserver: null,

			// handler functions set by render()
			cancelDrawHandler: null,
			deleteSelectionHandler: null,
		};
	},
	on_created_ready: function() {
		let params;
		try {
			params = JSON.parse(this.annoset.Params);
		} catch(e) {}
		if(!params) {
			params = {};
		}
		if(!params.Mode) {
			params.Mode = 'box';
		}
		if(!params.Categories) {
			params.Categories = [];
			params.CategoriesStr = '';
		} else {
			params.CategoriesStr = params.Categories.join(',');
		}
		this.params = params;

		// call handlers on certain key presses
		this.setKeyupHandler((e) => {
			if(document.activeElement.tagName == 'INPUT') {
				return;
			}

			if(e.key === 'Escape' && this.cancelDrawHandler) {
				this.cancelDrawHandler();
			} else if(e.key === 'Delete' && this.deleteSelectionHandler) {
				this.deleteSelectionHandler();
			}
		});
	},
	destroyed: function() {
		this.setKeyupHandler(null);
		this.disconnectResizeObserver();
	},
	on_update: function() {
		this.shapes = [];
		for(let i = 0; i < this.numFrames; i++) {
			this.shapes.push([]);
		}
	},
	on_item_data: function(data) {
		if(data.length == 0) {
			return;
		}
		this.shapes = data.map((shapeList) => {
			return shapeList.map((shp) => this.decodeShape(shp));
		});

		// update if we already rendered image
		if(this.imageDims != null) {
			this.render();
		}
	},
	on_image_loaded: function() {
		Vue.nextTick(() => {
			this.render();
		});
	},
	getAnnotateData: function() {
		let data = this.shapes.map((shapeList) => {
			return shapeList.map((shape) => this.encodeShape(shape))
		});
		let metadata = {
			CanvasDims: [this.imageDims.Width, this.imageDims.Height],
			Categories: this.params.Categories,
		};
		return [data, metadata];
	},
	methods: {
		disconnectResizeObserver: function() {
			if(this.resizeObserver) {
				this.resizeObserver.disconnect();
				this.resizeObserver = null;
			}
		},
		decodeShape: function(shape) {
			let shp = {};
			if(this.dataType === 'shape') {
				shp.Type = shape.Type;
				shp.Points = shape.Points;
			} else if(this.dataType === 'detection') {
				shp.Type = 'box';
				shp.Points = [[shape.Left, shape.Top], [shape.Right, shape.Bottom]];
			}
			shp.Category = (shape.Category) ? shape.Category : '';
			shp.TrackID = (shape.TrackID) ? shape.TrackID : '';
			return shp;
		},
		encodeShape: function(shape) {
			let shp = {};
			if(this.dataType === 'shape') {
				shp.Type = shape.Type;
				shp.Points = shape.Points;
			} else if(this.dataType === 'detection') {
				shp.Left = shape.Points[0][0];
				shp.Top = shape.Points[0][1];
				shp.Right = shape.Points[1][0];
