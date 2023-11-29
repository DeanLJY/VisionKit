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
				shp.Bottom = shape.Points[1][1];
			}
			if(shape.Category !== '') {
				shp.Category = shape.Category;
			}
			if(shape.TrackID !== '') {
				shp.TrackID = parseInt(shape.TrackID);
			}
			return shp;
		},
		updateCategories: function() {
			if(this.params.CategoriesStr == '') {
				this.params.Categories = [];
			} else {
				this.params.Categories = this.params.CategoriesStr.split(',');
			}
		},
		saveParams: function() {
			let request = {
				Params: JSON.stringify({
					Mode: this.params.Mode,
					Categories: this.params.Categories,
				}),
			};
			utils.request(this, 'POST', '/annotate-datasets/'+this.annoset.ID, JSON.stringify(request));
		},
		setKeyupHandler: function(handler) {
			if(this.keyupHandler != null) {
				this.$parent.$off('keyup', this.keyupHandler);
				this.keyupHandler = null;
			}
			if(handler != null) {
				this.keyupHandler = handler;
				this.$parent.$on('keyup', this.keyupHandler);
			}
		},
		render: function() {
			let stage = new Konva.Stage({
				container: this.$refs.layer,
				width: this.imageDims.Width,
				height: this.imageDims.Height,
			});
			let layer = new Konva.Layer();
			let resizeLayer = null;
			let destroyResizeLayer = () => {
				if(resizeLayer) {
					resizeLayer.destroy();
					resizeLayer = null;
				}
			};

			// we want annotations to be stored in coordinates based on image natural width/height
			// but in the UI, image could be stretched to different width/height
			// so here we need to stretch the stage in the same way
			let getScale = () => {
				return Math.min(
					this.$refs.image.width / this.imageDims.Width,
					this.$refs.image.height / this.imageDims.Height,
				);
			};
			let rescaleLayer = () => {
				if(!this.$refs.layer || !this.$refs.image) {
					return;
				}
				let scale = getScale();
				stage.width(parseInt(scale*this.imageDims.Width));
				stage.height(parseInt(scale*this.imageDims.Height));
				layer.scaleX(scale);
				layer.scaleY(scale);
				layer.draw();
				if(resizeLayer) {
					resizeLayer.scaleX(scale);
					resizeLayer.scaleY(scale);
					resizeLayer.draw();
				}
			};
			this.disconnectResizeObserver();
			this.resizeObserver = new ResizeObserver(rescaleLayer);
			this.resizeObserver.observe(this.$refs.image);
			rescaleLayer();
			let getPointerPosition = () => {
				let transform = layer.getAbsoluteTransform().copy();
				transform.invert();
				let pos = stage.getPointerPosition();
				return transform.point(pos);
			};

			let konvaShapes = [];
			// curShape is set if we are currently drawing a new shape
			let curShape = null;

			let resetColors = () => {
				konvaShapes.forEach((kshp, idx) => {
					if(this.selectedIdx === idx) {
						kshp.stroke('orange');
					} else {
						kshp.strok