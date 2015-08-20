package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

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
			// Parameters
		},
	}

	return resources
}

var swagify *bool

func init() {
	swagify = flag.Bool("swagify", false, "create a swag file for this server")
}

func main() {
	flag.Parse()
	resources := newResources()
	sg := snowflake.SG{}
	sg.Logger = log.New(os.Stdout, "", log.Lshortfile)

	gOptions := snowflake.GlobalOptions{
		SG:              sg,
		Port:            "2020",
		HealthcheckPort: "2023",
	}

	if *swagify {
		snowflake.Swagify(resources)
		return
	}

	snowflake.Run(resources, &gOptions)
}

func handler(sg snowflake.SG, w http.ResponseWriter, r *http.Request) {
	sg.Logger.Printf("handler log line")
	fmt.Fprintf(w, "Handled")
}
