import PytorchTrainGeneric from './pytorch_train-generic.js';
import SelectInputSize from './select-input-size.vue';

export default PytorchTrainGeneric({
	components: {
		'select-input-size': SelectInputSize,
	},
	disabled: ['model', 'dataset'],
	created: function() {
		if(!('Mode' in this.params)) this.$set(this.params, 'Mode', 'yolov3');
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
		<label class="col-sm-4 col-form-label">Mode</label>
		<div class="col-sm-8">
			<select v-model="params.Mode" class="form-select">
				<option value="yolov3">YOLOv3</option>
				<option value="yolov3-spp">YOLOv3-SPP</option>
				<option value="yolov3-tiny">YOLOv3-Tiny</option>
			</select>
			<small class="form-text text-muted">
				YOLOv3 and YOLOv3-SPP are large models providing high accuracy (YOLOv3-SPP may provide slightly higher accuarcy).
				YOLOv3-Tiny is a small model that is fast but provides lower accuracy.
			</small>
		</div>
	</div>
	<div class="row mb-2">
		<label class="col-sm-4 col-form-label">Validati