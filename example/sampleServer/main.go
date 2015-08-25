package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/trotha01/snowflake"
)

type handle struct {
	LogChan    chan string
	MetricChan chan string
	AlertChan  chan string
}

func newResources(logChan, metricChan, alertChan chan string) snowflake.Resources {

	localOptions := snowflake.ResourceOptions{
		Timeout:   "60",
		RateLimit: "150",
	}

	h := handle{
		LogChan:    logChan,
		MetricChan: metricChan,
		AlertChan:  alertChan,
	}

	resources := snowflake.Resources{
		{
			Path:    "/",
			Get:     h,
			Options: localOptions,
		},
		{
			Path: "/v3/endpoint",
			Post: h,
			Parameters: []snowflake.Parameter{
				snowflake.Parameter{
					Name:        "param",
					Description: "a parameter for this endpoint",
					Type:        "int",
				},
			},
		},
		{
			Path: "/seth/endpoint",
			Get:  h,
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
	logChan := make(chan string)
	metricChan := make(chan string)
	alertChan := make(chan string)
	resources := newResources(logChan, metricChan, alertChan)

	if *swagify {
		snowflake.Swagify(resources)
		return
	}

	gOptions := snowflake.GlobalOptions{
		LogChan:         logChan,
		Port:            "2020",
		HealthcheckPort: "2023",
	}

	go logListener(logChan)
	go metricListener(metricChan)
	go alertListener(alertChan)

	snowflake.Run(resources, &gOptions)
}

func logListener(logChan chan string) {
	for {
		message := <-logChan
		log.Printf(message)
	}
}

func metricListener(metricChan chan string) {
	for {
		message := <-metricChan
		log.Println("metric: ", message)
	}
}

func alertListener(alertChan chan string) {
	for {
		message := <-alertChan
		log.Println("alert: ", message)
	}
}

func (h handle) Handle(id int, w http.ResponseWriter, r *http.Request) {
	h.LogChan <- "handler log line through chan"
	h.MetricChan <- "Add this metric"
	h.AlertChan <- "ALERT!"
	fmt.Fprintf(w, "Handled")
}
