package main

import (
	"fmt"
	"net/http"

	"github.com/trotha01/snowflake"
)

func newResources() snowflake.Resources {

	someOptions := snowflake.ResourceOptions{
		Timeout:   "60",
		RateLimit: "150",
	}

	resources := snowflake.Resources{
		{
			Path:    "/",
			Get:     handler,
			Options: someOptions,
		},
		{
			Path: "/v3/endpoint",
			Post: handler,
		},
	}

	return resources
}

func main() {
	resources := newResources()

	gOptions := snowflake.GlobalOptions{Port: "2020", HealthcheckPort: "2023"}

	snowflake.Run(resources, &gOptions)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Handled")
}
