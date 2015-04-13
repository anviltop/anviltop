package api

import (
	"fmt"
	"github.com/anviltop/anviltop/models"
	"net/http"
)

func Handler(res http.ResponseWriter, req *http.Request) {

	fmt.Print("In Job Definition Handler")
	DumpReq(req)

	jobDefinition := models.ExampleContainerJob()

	WriteObjectResponse(res, jobDefinition)
}
