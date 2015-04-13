package models

type JobOutput struct {

	// Identifier for the job
	JobId string

	// Output - store as byte array
	Output string
}

func (jobOutput JobOutput) RedisKeyPrefix() string {
	result := "JobOutput"
	return result
}

func (jobOutput JobOutput) Id() string {
	return jobOutput.JobId
}
