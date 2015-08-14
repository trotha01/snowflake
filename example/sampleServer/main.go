package main

import (
	"fmt"
	"net/http"

	"github.com/trotha01/snowflake"
)

func newResources() snowflake.Resources {
	myOptions := snowflake.ResourceOptions{
		Timeout:   "60",
		RateLimit: "150",
	}

	resources := snowflake.Resources{
		{
			Path:    "/",
			Get:     handler,
			Options: myOptions,
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
	// return c.String(http.StatusOK, "Worked!")
	fmt.Fprintf(w, "Hello")
}
