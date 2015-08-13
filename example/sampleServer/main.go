package main

import (
	"fmt"
	"net/http"

	"github.com/trotha01/snowflake"
)

var Resources snowflake.Resources

func main() {
	myOptions := snowflake.ResourceOptions{
		Timeout:   "60",
		RateLimit: "150",
	}

	Resources := snowflake.Resources{
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

	gOptions := snowflake.GlobalOptions{Port: "2020", HealthcheckPort: "2023"}
	snowflake.Run(Resources, gOptions)
}

func handler(w http.ResponseWriter, r *http.Request) {
	// return c.String(http.StatusOK, "Worked!")
	fmt.Fprintf(w, "Hello")
}
