package snowflake

import (
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/gommon/color"
)

type Resources []Resource

// type Handler func(http.ResponseWriter, *http.Request)
type Handler http.HandlerFunc

type SG struct {
	Logger log.Logger
	// StatsD client
	// APId client
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
	// e.Use(mWare)
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

func Logger() echo.MiddlewareFunc {
	return func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			// h(c.Response().Writer(), c.Request)
			// return nil

			start := time.Now()
			if err := h(c); err != nil {
				c.Error(err)
			}
			stop := time.Now()
			method := c.Request().Method
			path := c.Request().URL.Path
			if path == "" {
				path = "/"
			}
			size := c.Response().Size()

			n := c.Response().Status()
			code := color.Green(n)
			switch {
			case n >= 500:
				code = color.Red(n)
			case n >= 400:
				code = color.Yellow(n)
			case n >= 300:
				code = color.Cyan(n)
			}

			log.Printf("%s %s %s %s %d", method, path, code, stop.Sub(start), size)
			return nil
		}
	}
}
