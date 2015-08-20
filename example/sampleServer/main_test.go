package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/trotha01/snowflake"
)

var testClient http.Client

func init() {
	testClient = snowflake.MockRun(newResources(), nil)
}

func TestRoot(t *testing.T) {
	res, err := testClient.Get("/")
	if err != nil {
		t.Fatal("Error with get request to /. Error: " + err.Error())
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Fatalf("Error reading request body - %s", err.Error())
	}
	fmt.Printf("body: %s\n", body)
	// do verifications
}

func TestEndpoint(t *testing.T) {
	res, err := testClient.Post("/v3/endpoint", "", nil)
	if err != nil {
		t.Fatal("Error with get request to /v3/endpoint. Error: " + err.Error())
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Fatalf("Error reading request body - %s", err.Error())
	}
	fmt.Printf("body: %s\n", body)
	// do verifications
}
