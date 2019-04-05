package engine

import (
	"time"
)

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
			i := int64(sr.ExecutionTime / time.Millisecond)
			if i == 0 {
				sr.ExecutionTime, _ = time.ParseDuration("1ms")
			}

			r.Steps.Results = append(r.Steps.Results, sr)
			r.Steps.TotalExecutionTime = r.Steps.TotalExecutionTime + sr.ExecutionTime

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
