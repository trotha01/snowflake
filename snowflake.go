package snowflake

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

type Resources []Resource

// A object for logging and clients
type SG struct {
	Logger *log.Logger
	// Statsd client
}

// type Handler http.HandlerFunc
type Handler func(SG, http.ResponseWriter, *http.Request)

type Resource struct {
	Path       string
	Get        Handler
	Post       Handler
	Delete     Handler
	Put        Handler
	Patch      Handler
	Options    ResourceOptions
	Headers    []string
	Parameters []string
}

type ResourceOptions struct {
	Timeout   string
	RateLimit string
}

type GlobalOptions struct {
	Timeout         string
	RateLimit       string
	Port            string
	HealthcheckPort string
	SG              SG
}

// Run starts the http server
func Run(resources Resources, options *GlobalOptions) {
	e := echo.New()
	for _, resource := range resources {
		if resource.Get != nil {
			e.Get(resource.Path, options.SG.middlewareHandler(resource.Get))
		}
		if resource.Post != nil {
			e.Post(resource.Path, options.SG.middlewareHandler(resource.Post))
		}
		if resource.Delete != nil {
			e.Delete(resource.Path, options.SG.middlewareHandler(resource.Delete))
		}
	}

	options.SG.Logger.Printf("Running on port " + options.Port)
	go e.Run(":" + options.Port)

	eHealth := echo.New()
	eHealth.Get("/healthcheck", healthHandler)
	options.SG.Logger.Printf("Running Healthcheck on port " + options.HealthcheckPort)
	eHealth.Run(":" + options.HealthcheckPort)
}

func healthHandler(c *echo.Context) error {
	c.String(200, "all good here")
	return nil
}

// turns a http handler into a echo handler
func (sg SG) middlewareHandler(handler Handler) echo.HandlerFunc {
	return func(c *echo.Context) error {
		sg.Logger.SetPrefix(strconv.Itoa(rand.Int()) + " ")
		handler(sg, c.Response().Writer(), c.Request())
		return nil
	}
}

func Swagify(resources Resources) {
	header := `
swagger: '2.0'
info:
  title: Uber API
  description: Move your app forward with the Uber API
  version: "1.0.0"

host: api.server.com
schemes:
  - https
produces:
  - application/jston
paths:
`
	fmt.Printf(header)

	for _, resource := range resources {
		if resource.Get != nil {
			fmt.Printf("  %s:\n", resource.Path)
			fmt.Printf("    get:\n")
		}
		if resource.Post != nil {
			fmt.Printf("  %s:\n", resource.Path)
			fmt.Printf("    post:\n")
		}

	}

}
