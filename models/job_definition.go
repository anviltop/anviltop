package models

import (
	"bytes"
)

type JobDefinition struct {

	// Identifier for the models id
	JobDefinitionId string

	// A set of options
	Options map[string]string

	// The {organization}/{repository}
	Image string

	// The specific container version to pull
	ImageTag string

	// Environment variables
	Environment map[string]string

	// The command to run
	Command string

	// Arguments to the command
	CommandArgs map[string]string
}

func (jobDefinition JobDefinition) RedisKeyPrefix() string {
	result := "JobDefinition"
	return result
}

func (jobDefinition JobDefinition) Id() string {
	return jobDefinition.JobDefinitionId
}

func (jobDefinition *JobDefinition) BashString() string {

	var stringBuffer bytes.Buffer

	// Docker run
	stringBuffer.WriteString("docker run")

	// Run Options
	if jobDefinition.Options != nil {
		for optionName, optionValue := range jobDefinition.Options {
			stringBuffer.WriteString(" " + optionName + " " + optionValue)
		}
	}

	// Environment
	if jobDefinition.Environment != nil {
		for envVariable, envValue := range jobDefinition.Environment {
			stringBuffer.WriteString(" -e \"" + envVariable + "=" + envValue + "\"")
		}
	}

	// Container Image
	stringBuffer.WriteString(" " + jobDefinition.Image)
	if jobDefinition.ImageTag != "" {
		stringBuffer.WriteString(":" + jobDefinition.ImageTag)
	}

	// Command
	if jobDefinition.Command != "" {
		stringBuffer.WriteString(" " + jobDefinition.Command)
	}

	result := stringBuffer.String()
	return result
}

func ExampleContainerJob() JobDefinition {
	result := JobDefinition{
		Image:       "docker_example_crontainers/example_1_failure",
		ImageTag:    "",
		Command:     "",
		Options:     nil,
		Environment: nil,
	}

	result.Options = make(map[string]string)
	result.Options["--name"] = "test"
	result.Options["--rm"] = ""
	result.Environment = make(map[string]string)
	result.Environment["TEST"] = "TESTVALUE"

	return result
}
