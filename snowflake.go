package snowflake

import (
	"log"

	"github.com/labstack/echo"
)

type Resources []Resource

type Resource struct {
	Path       string
	Get        func()
	Post       func()
	Del        func()
	Options    Options
	Headers    []string
	Parameters []string
}

type Options struct {
	Timeout   string
	RateLimit string
}

func Run(resources Resources) {
	// run http server
	log.Printf("%+v", resources)

	e := echo.New()
	for _, resource := range resources {
		if resource.Get != nil {
			e.Get(resource.Path, resource.Get)
		}
	}
	e.Run(":2020")
}
