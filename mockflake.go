package snowflake

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
)

type testRoundTripper struct {
	resources Resources
	recorder  *httptest.ResponseRecorder
	sg        SG
}

func MockRun(r Resources, options *GlobalOptions) http.Client {

	sg := SG{
		Logger: log.New(os.Stdout, "", log.Lshortfile),
	}

	roundTripper := testRoundTripper{
		resources: r,
		recorder:  httptest.NewRecorder(),
		sg:        sg,
	}

	testClient := http.Client{
		Transport: roundTripper,
	}

	return testClient
}

func (trt testRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	w := trt.recorder

	for _, resource := range trt.resources {
		fmt.Printf("resource path: %q\n", resource.Path)
		fmt.Printf("request path: %q\n", request.URL.Path)
		if resource.Path == request.URL.Path {
			if request.Method == "GET" {
				resource.Get(trt.sg, w, request)
				return respFromRecorder(w), nil
			}
			if request.Method == "POST" {
				resource.Post(trt.sg, w, request)
				return respFromRecorder(w), nil
			}
		}
	}

	return nil, errors.New("path not found")
}

// respFromRecorder builds an http response from a httptest recorder
func respFromRecorder(w *httptest.ResponseRecorder) *http.Response {
	resp := http.Response{}
	resp.StatusCode = w.Code
	resp.Header = w.Header()
	// TODO: fill in the rest of response

	b := w.Body.Bytes()
	resp.Body = ioutil.NopCloser(bytes.NewReader(b))
	return &resp
}
