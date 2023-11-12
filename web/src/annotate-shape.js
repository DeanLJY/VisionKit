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
			selectedId