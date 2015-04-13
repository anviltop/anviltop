package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"

	"github.com/anviltop/anviltop/controllers/api"
)

func echoHandler(res http.ResponseWriter, req *http.Request) {

	b, err := json.Marshal(req)
	if err != nil {
		fmt.Println("error:", err)
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}

	fmt.Print("Req: ", string(b), "\n")
	fmt.Print("Body: ", string(body), "\n")
	fmt.Fprintf(res, "OK")
}

func main() {

	// Create the redis client - Done by the var block above

	// Define routes and handlers
	router := mux.NewRouter()
	router.HandleFunc("/api/job/reporting/{job_id}/start", api.JobStartHandler)
	router.HandleFunc("/api/job/reporting/{job_id}/output", api.JobOutputHandler)
	router.HandleFunc("/api/job/reporting/{job_id}/exit_status", api.JobExitStatusHandler)

	http.Handle("/", router)
	http.ListenAndServe(":8082", nil)
	fmt.Printf("Start web server")
}
