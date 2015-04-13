package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func LogReq(req *http.Request) {
	fmt.Println(req.Method + " " + req.URL.Path)
}

func DumpReq(req *http.Request) {
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
}

func WriteObjectResponse(res http.ResponseWriter, object interface{}) {
	// serialize as JSON
	raw, err := json.Marshal(object)
	if err != nil {
		fmt.Print("ERROR SERIALIZING!")
		// return HTTP 500
		res.Write([]byte("ERROR"))
	} else {
		res.Write(raw)
	}
}
