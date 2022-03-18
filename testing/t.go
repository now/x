package testing

// T is used by test helpers to report fatal states and logging.
type T interface {
	// Fail the current test.
	//
	// Execution of the test continues, but the test run is marked as having failed.
	Fail()

	// Fatal state has been entered and execution should be stopped after reporting arguments.
	Fatal(...interface{})

	// Helper marks the calling function as a test helper function.
	Helper()

	// Log arguments using fmt.Sprintln(arguments...).
	//
	// The log is displayed for failing tests or if testing.Verbose() is true.
	Log(...interface{})
}
