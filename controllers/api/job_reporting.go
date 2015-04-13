package api

import (
	"fmt"
	"github.com/anviltop/anviltop/models"
	"github.com/anviltop/anviltop/util/redis"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strconv"
)

func JobStartHandler(res http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	LogReq(req)
	vars := mux.Vars(req)

	var job models.Job

	// Get job id from request
	jobId, exists := vars["job_id"]
	if !exists || jobId == "0" {

		fmt.Println("No job id, create new job")

		// Create a new job if it doesn't exist; save it
		job = models.Job{}
		job.JobId = models.NextId("Job") // get new job id
		job.Hostname = req.FormValue("hostname")
		job.Cmd = req.FormValue("cmd")

	} else if jobId != "" && jobId != "0" {

		// Get the job from database, fake for now
		err := models.Get(&job, jobId)
		if err != nil {
			fmt.Println("Error Getting Model: ", err)
		}
	}

	// Mark this job as running
	job.Status = models.JOB_STATUS_RUNNING

	// Save the object to redis
	models.Set(&job)

	// Add this job to running job list
	redis.SAdd(redis.RUNNING_JOBS_KEY, jobId)

	// Respond with a job object
	WriteObjectResponse(res, job)
}

func JobExitStatusHandler(res http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	LogReq(req)
	vars := mux.Vars(req)

	var job models.Job

	// Get job id from request
	jobId, exists := vars["job_id"]
	if !exists || jobId == "0" {
		fmt.Println("No job id passed!")
		return
	}

	// Get the job from redis
	err := models.Get(&job, jobId)
	if err != nil {
		fmt.Println("Error Getting Model: ", err)
	}

	// Update the job state with completed and exit status
	job.Status = models.JOB_STATUS_COMPLETED
	exitStatus, _ := strconv.Atoi(req.FormValue("exit_status"))
	job.ExitStatus = exitStatus

	// Save the job state
	models.Set(&job)

	// Remove this job from the running job list
	redis.SRem(redis.RUNNING_JOBS_KEY, jobId)

	//TODO: Save the JobOutput to a file?

	// Respond with the job object
	WriteObjectResponse(res, job)
}

func JobOutputHandler(res http.ResponseWriter, req *http.Request) {
	LogReq(req)
	vars := mux.Vars(req)

	// Get job id from request
	jobId, exists := vars["job_id"]
	if !exists || jobId == "0" {
		fmt.Println("No job id passed!")
		return
	}

	jobOutput := models.JobOutput{JobId: jobId}

	// Get the job output from redis
	err := models.Get(&jobOutput, jobId)
	if err != nil {
		if err.Error() == "redigo: nil returned" {
			fmt.Println("NIL RETURNED ERROR!")
		}
		fmt.Println("Error Getting Model: ", err)
	}

	// Read the body and append it to the job output
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}
	bodyStringA := string(body)

	// Append string - naive approach
	jobOutput.Output += bodyStringA

	fmt.Println("JobOutput: ", jobOutput.Output)

	models.Set(&jobOutput)

	fmt.Fprintf(res, "OK")
}
