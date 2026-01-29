package gocron

// Must panics if err is non-nil and otherwise returns the job.
func Must(j Job, err error) Job {
	if err != nil {
		panic(err)
	}

	return j
}
