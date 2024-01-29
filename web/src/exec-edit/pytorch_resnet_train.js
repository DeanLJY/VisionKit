import PytorchTrainGeneric from './pytorch_train-generic.js';
import SelectInputSize from './select-input-size.vue';

export default PytorchTrainGeneric({
	components: {
		'select-input-size': SelectInputSize,
	},
	disabled: ['model', 'dataset'],
	created: function() {
		if(!('Mode' in this.params)) this.$set(this.params, 'Mode', 'yolov3');
		if(!('NumClasses' in this.params)) this.$set(this.params, 'NumClasses', 0);
		if(!('ValPercent' in this.params)) this.$set(this.params, 'ValPercent', 20);
		if(!('Resize' in this.params)) {
			this.$set(this.params, 'Resize', {
				Mode: 'fixed',
				MaxDimension: 640,
				Width: 224,
				Height: 224,
				Multiple: 1,
			});
		}
	},
	basicTemplate: `
<div class="small-container">
	<div class="row mb-2">
		<label class="col-sm-4 col-form-label">Mode</label>
		<div class="col-sm-8">
			<select v-model="params.Mode" class="form-select">
				<option value="resnet18">Resnet18</option>
				<option value="resnet34">Resnet34</option>
				<option value="resnet50">Resnet50</option>
				<option value="resnet101">Resnet101</option>
				<option value="resnet152">Re