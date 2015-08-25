package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

// The purpose of this tool is to replay requests

// request is to be retrieved from a log line
// the log line "Message" should be "request"

const nonRequestLogErr = "not a request log"
const Host = "localhost"
const Port = "2020"

// TODO: rename
type logStruct struct {
	Message string `json:"message"`
}

type request struct {
	Header     http.Header `json:"header"`
	Body       string      `json:"body"`
	Host       string      `json:"host"`
	Path       string      `json:"path"`
	Message    string      `json:"message"`
	Method     string      `json:"method"`
	RemoteAddr string      `json:"remoteAddr"`
}

func main() {
	logLinesChan := make(chan []byte)
	var wg sync.WaitGroup
	wg.Add(1)
	go logLines(logLinesChan, &wg)
	go handleLogLines(logLinesChan)
	wg.Wait()
}

func handleLogLines(logLinesChan chan []byte) {
	for logLine := range logLinesChan {
		// ignore empty lines
		if bytes.Equal(bytes.Trim(logLine, " \n"), []byte("")) {
			continue
		}
		logRequest, err := requestFromLogLine(logLine)
		// ignore logs that aren't a request
		if err != nil && err.Error() == nonRequestLogErr {
			continue
		}
		if err != nil {
			log.Printf("Error parsing log line: %s, error: %s\n", logLine, err.Error())
			continue
		}

		err = doRequest(logRequest)
		if err != nil {
			log.Printf("Error executing log request: %s", err.Error())
		}
	}

}

func doRequest(logRequest request) error {
	host := fmt.Sprintf("http://%s:%s", Host, Port)
	body := strings.NewReader(logRequest.Body)

	req, err := http.NewRequest(logRequest.Method, host+logRequest.Path, body)
	if err != nil {
		return err
	}
	req.Header = logRequest.Header

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

// Sends log lines, line by line, to the channel
func logLines(logLinesChan chan []byte, wg *sync.WaitGroup) {
	reader := bufio.NewReader(os.Stdin)
	for {
		l, err := reader.ReadBytes('\n')
		if err != nil {
			if err != io.EOF {
				log.Fatalf("Error reading stdin: %s", err.Error())
			}
			break
		}
		logLinesChan <- l
	}

	wg.Done()
}

func requestFromLogLine(logLine []byte) (request, error) {
	var r request
	var l logStruct
	// get log starting after the "{"
	i := bytes.Index(logLine, []byte("{"))
	if i == -1 {
		return r, errors.New("Not a properly formed log line, missing json struct")
	}
	err := json.Unmarshal(logLine[i:], &l)
	if err != nil {
		return r, err
	}
	if l.Message != "request" {
		return r, fmt.Errorf(nonRequestLogErr)
	}

	err = json.Unmarshal(logLine[i:], &r)
	if err != nil {
		return r, err
	}
	return r, nil
}
