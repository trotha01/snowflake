package main

import "github.com/trotha01/snowflake"

func main() {
	myOptions := snowflake.Options{
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

	snowflake.Run(resources)
}

func handler() {
}
