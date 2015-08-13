package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

// type testResponseWriter http.ResponseWriter

type testTransport struct{}

type testResponseWriter struct {
	header            http.Header
	data              []byte
	statusCode        int
	writeHeaderCalled bool
	writeCalled       bool
}

var testClient *http.Client

// RoundTrip must be safe for concurrent use by multiple goroutines per the http package documentation
func (testTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	var response http.Response
	var responseWriter = testResponseWriter{}

	for _, resource := range Resources {
		if resource.Path == request.RequestURI {
			if request.Method == "GET" {
				resource.Get(responseWriter, request)
				return &response, nil
			}
		}
	}

	return nil, errors.New("resource does not exist - " + request.RequestURI)
}

func init() {
	transport := testTransport{}
	testClient = &http.Client{
		Transport: transport,
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		for _, resource := range Resources {
			if resource.Path == request.RequestURI {
				if request.Method == "GET" {
					resource.Get(w, request)
				}
			}
		}

	}))
	defer ts.Close()

}

func TestRoot(t *testing.T) {
	testClient.Get("/")
	// http.Get("/")
	// do verifications
}

func TestEndpoint(t *testing.T) {
	// http.Get("/v3/endpoint")
	// do verifications
}

func newTestResponseWriter() http.ResponseWriter {
	return testResponseWriter{
		header:     http.Header{},
		data:       []byte{},
		statusCode: -1,
	}
}

func (trw testResponseWriter) Header() http.Header {
	return trw.header
}

func (trw testResponseWriter) Write(data []byte) (int, error) {
	if trw.statusCode == -1 {
		trw.statusCode = http.StatusOK
	}
	trw.data = data
	trw.writeCalled = true
	return 0, nil
}

func (trw testResponseWriter) WriteHeader(statusCode int) {
	// TODO: should I do this? Test with actual http server
	// if write has already been called, do nothing
	// if trw.writeCalled {
	// 	return
	// }

	trw.statusCode = statusCode
	trw.writeHeaderCalled = true
}
