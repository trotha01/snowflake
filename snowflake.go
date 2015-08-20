package snowflake

import (
	"log"
	"net/http"

	"github.com/labstack/echo"
)

type Resources []Resource

type Handler http.HandlerFunc

type SG struct {
	Logger log.Logger
	// Statsd client
}

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
}

// Run starts the http server
func Run(resources Resources, options *GlobalOptions) {
	e := echo.New()
	for _, resource := range resources {
		if resource.Get != nil {
			e.Get(resource.Path, middlewareHandler(resource.Get))
			// e.Get(resource.Path, http.HandlerFunc(resource.Get))
		}
		if resource.Post != nil {
			e.Post(resource.Path, middlewareHandler(resource.Post))
			// e.Post(resource.Path, http.HandlerFunc(resource.Post))
		}
		if resource.Delete != nil {
			e.Delete(resource.Path, middlewareHandler(resource.Delete))
			// e.Delete(resource.Path, http.HandlerFunc(resource.Delete))
		}
	}

	log.Printf("Running on port " + options.Port)
	go e.Run(":" + options.Port)

	eHealth := echo.New()
	eHealth.Get("/healthcheck", healthHandler)
	log.Printf("Running Healthcheck on port " + options.HealthcheckPort)
	eHealth.Run(":" + options.HealthcheckPort)
}

func healthHandler(c *echo.Context) error {
	c.String(200, "all good here")
	return nil
}

// turns a http handler into a echo handler
func middlewareHandler(handler Handler) echo.HandlerFunc {
	return func(c *echo.Context) error {
		handler(c.Response().Writer(), c.Request())
		return nil
	}
}
