package snowflake

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"

	"github.com/labstack/echo"
)

type Resources []Resource

// The int will be the request id
type Handler interface {
	Handle(int, http.ResponseWriter, *http.Request)
}

type Parameters []Parameter

type Parameter struct {
	Name        string
	Description string
	Type        string
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
	Parameters Parameters
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
	LogChan         chan string
}

// Run starts the http server
func Run(resources Resources, options *GlobalOptions) {
	e := echo.New()
	for _, resource := range resources {
		if resource.Get != nil {
			e.Get(resource.Path, middlewareHandler(resource.Get, options.LogChan))
		}
		if resource.Post != nil {
			e.Post(resource.Path, middlewareHandler(resource.Post, options.LogChan))
		}
		if resource.Delete != nil {
			e.Delete(resource.Path, middlewareHandler(resource.Delete, options.LogChan))
		}
	}

	options.LogChan <- fmt.Sprintf("Running on port " + options.Port)
	go e.Run(":" + options.Port)

	eHealth := echo.New()
	eHealth.Get("/healthcheck", healthHandler)
	options.LogChan <- fmt.Sprintf("Running Healthcheck on port " + options.HealthcheckPort)
	eHealth.Run(":" + options.HealthcheckPort)
}

func healthHandler(c *echo.Context) error {
	c.String(200, "all good here")
	return nil
}

// turns a http handler into a echo handler
func middlewareHandler(handler Handler, logChan chan string) echo.HandlerFunc {
	return func(c *echo.Context) error {
		r := c.Request()

		buf := bytes.NewBuffer(nil)
		_, err := io.Copy(buf, r.Body)
		r.Body.Close()

		h, err := json.Marshal(r.Header)
		if err != nil {
			logChan <- "err marshelling request header. Error: " + err.Error()
			return errors.New("Internal Server Error marshalling request header")
		}
		logChan <- fmt.Sprintf(`{"message":"request", "remoteAddr": "%s", "method": "%s", "host": "%s", "path": "%s", "body": "%s", "header": %s}`, r.RemoteAddr, r.Method, r.URL.Host, r.URL.Path, buf.Bytes(), string(h)) // strconv.Quote(string(h)))

		if err != nil {
			logChan <- "err reading request body. Error: " + err.Error()
			return errors.New("Internal Server Error reading request body")
		}

		handler.Handle(rand.Int(), c.Response().Writer(), r)
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
			fmt.Printf("      parameters:\n")
			for _, parameter := range resource.Parameters {
				fmt.Printf("      name:%s\n", parameter.Name)
				fmt.Printf("      type:%s\n", parameter.Type)
				fmt.Printf("      description:%T\n", parameter.Description)
			}
		}
		if resource.Post != nil {
			fmt.Printf("  %s:\n", resource.Path)
			fmt.Printf("    post:\n")
			fmt.Printf("      parameters:\n")
			for _, parameter := range resource.Parameters {
				fmt.Printf("        name:%s\n", parameter.Name)
				fmt.Printf("        type:%s\n", parameter.Type)
				fmt.Printf("        description:%T\n", parameter.Description)
			}
		}
	}
}
