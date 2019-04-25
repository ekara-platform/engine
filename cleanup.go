package engine

// Cleanup represents a cleanup method to rollback what has been done by a step
type Cleanup func(lC LaunchContext) error

//NoCleanUpRequired is an predefined  cleanup which can be used to indicate that
// no cleanup is required
func NoCleanUpRequired(lC LaunchContext) error {
	// Do nothing and it's okay...
	// This is just an explicit empty implementation to clearly materialize that no cleanup is required
	return nil
}

//cleanLaunched runs a slice of cleanup functions
func cleanLaunched(fs []Cleanup, lC LaunchContext) (e error) {
	for _, f := range fs {
		e := f(lC)
		if e != nil {
			return e
		}
	}
	return nil
}
