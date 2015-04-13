package main

import (
	"github.com/anviltop/anviltop/remote"
)

func main() {

	// Get the job definition for what is being asked to be run
	//jobInstance := models.ExampleContainerJob()
	//jobInstance := remote.GetJobDefinition("100")
	//fmt.Print("Run ", jobInstance.BashString(), "\n")

	// Make sure we have the container here if needed
	/*
		if jobInstance.RunsInContainer() == true {
			containerImage = jobInstance.Image
			containerImageTag = jobInstance.ImageTag
		}
	*/

	// Run the command, with reporting - no job definition present
	bashString := "/Users/lrajlich/Projects/lrajlich/docker-example-crontainers/example_1_failure/batch_job.sh"
	remote.RunShellCommand(bashString, "0")
}
