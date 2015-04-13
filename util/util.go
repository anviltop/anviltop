package util

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func VarDump(object interface{}) string {
	raw, err := json.Marshal(object)
	if err != nil {
		// do nothing
	}
	result := string(raw)
	return result
}

func ParseResponse(response *http.Response, object interface{}) {

	// Get the body of the response and deserialize
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	//func (c *Command) LoadFromJSON(jsonStr string) error {
	//var data = &c
	err = json.Unmarshal(body, object)
}

func NoErr(object interface{}, err error) interface{} {
	return object
}
