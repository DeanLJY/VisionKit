import PytorchTrainGeneric from './pytorch_train-generic.js';
import SelectInputSize from './select-input-size.vue';

export default PytorchTrainGeneric({
	components: {
		'select-input-size': SelectInputSize,
	},
	disabled: ['model', 'dataset'],
	created: function() {
		if(!('Mode' in this.params)) this.$set(this.params, 'Mode', 'x');
		if(!('ValPercent' in this.params)) this.$set(this.params, 'ValPercent', 20);
		if(!('Resize' in this.params)) {
			this.$set(this.params, 'Resize', {
				Mode: 'scale-down',
				MaxDimension: 640,
				Width: 416,
				Height: 416,
				Multiple: 32,
			});
		}
	},
	basicTemplate: `
<div class="small-container">
	<div class="row mb-2">
		<label class="col-sm-4 col-form-label">Mode</lab