package models

import ()

const (
	JOB_STATUS_RUNNING   = "running"
	JOB_STATUS_COMPLETED = "complete"
)

type Job struct {

	// Identifier for the job
	JobId string

	// Identifier for the models id
	//JobDefinition JobDefinition
	JobDefinitionId string

	// If running via command line,
	Cmd string

	// Which host was this running on
	Hostname string

	// Which worker was this running on, if any (identify by host:port)
	Worker string

	// Value from a const above
	Status string

	ExitStatus int

	jobDefinition *JobDefinition
}

func (job Job) RedisKeyPrefix() string {
	result := "Job"
	return result
}

func (job Job) Id() string {
	return job.JobId
}

func (job Job) JobDefinition() *JobDefinition {
	if job.jobDefinition == nil {
		jobDefinition := JobDefinition{}
		Get(jobDefinition, job.JobDefinitionId)
	}
	return job.jobDefinition
}
