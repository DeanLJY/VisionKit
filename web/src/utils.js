function request(comp, method, endpoint, params, successFunc, completeFunc, opts) {
	var args = {
		type: method,
		url: endpoint,
		error: (req, status, errorMsg) => {
			if(!comp) {
				return;
			}
			if(req && req.responseText) {
				errorMsg = req.responseText;
			}
			if(comp.setError) {
				comp.setError(errorMsg);
			} else if(comp.$globals.app) {
				comp.$globals.app.setError(errorMsg);
			}
		},
	};
	if(params) {
		args.data = params;
		if(typeof(args.data) === 'string') {
			args.processData = false;
		}
	}
	if(successFunc) {
		args.success = successFunc;
	}
	if(completeFunc) {
		args.complete = completeFunc;
	}
	if(opts) {
		if(opts.dataType) {
			args.dataType = opts.dataType;
		}
		if(opts.error) {
			// override the error handler we