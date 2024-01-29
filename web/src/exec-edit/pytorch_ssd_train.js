import PytorchTrainGeneric from './pytorch_train-generic.js';
export default PytorchTrainGeneric({
	disabled: ['model', 'dataset'],
	created: function() {
		if(!('Mode' in this.params)) this.$set(this.params, 'Mode', 'mb2-ssd-lite');
		if(!('ValPercent' in this.p