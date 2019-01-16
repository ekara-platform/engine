package engine

// launch runs a slice of step functions
//
// If one step in the slice returns an error then the launch process will stop and
// the cleanup will be invoked on all previously launched steps
func launch(fs []step, lC LaunchContext, rC *runtimeContext) ExecutionReport {

	r := &ExecutionReport{
		Context: lC,
	}

	cleanups := []Cleanup{}

	for _, f := range fs {
		ctx := f(lC, rC)
		for _, sr := range ctx.Results {
			r.Steps.Results = append(r.Steps.Results, sr)
			if sr.cleanUp != nil {
				cleanups = append(cleanups, sr.cleanUp)
			}

			e := sr.error
			if e != nil {
				cleanLaunched(cleanups, lC)
				r.Error = e
				return *r
			}
		}
	}
	return *r
}
