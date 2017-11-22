package server

import (
	"fmt"
	"net/url"
	"strings"
	"time"
)

// File struct definition
type File struct {
	FileName string
	Content  string
}

// Request struct definition
type Request struct {
	Method      string
	URL         *url.URL
	Httpversion string
	Headers     map[string]string
	Body        string
}

// Response struct definition
type Response struct {
	HTTPVersion string
	Status      string
	Error       string
	Headers     map[string]string
	Body        string
}

//converts response to string
func (response Response) ToString() (responseText string) {
	responseText = fmt.Sprintf("%s %s %s \r\n", response.HTTPVersion, response.Error, response.Status)
	response.Headers["Server"] = "COMP445/2.0 (Assignment)"
	now := time.Now()
	response.Headers["Last-Modified"] = now.Format("Mon Jan 2 15:04:05 MST 2006")
	response.Headers["Content-Length"] = fmt.Sprintf("%d", len(response.Body))
	for name, value := range response.Headers {
		responseText += fmt.Sprintf("%s: %s \r\n", name, value)
	}
	responseText += fmt.Sprintf("%s \r\n\r\n", response.Body)
	return
}

//Function to take buffer data and parse it into Request
func ParseRequest(buf []byte) (request Request) {
	var lines []string
	//Split each line
	lines = strings.Split(string(buf), "\r\n")

	// Split the first line with the request definition
	head := strings.Split(lines[0], " ")

	if len(head) > 2 {
		u, _ := url.Parse(head[1])

		headers := make(map[string]string)
		var body string
		isBodyData := false
		for i := 1; i < len(lines); i += 2 {
			if lines[i] != "" && !isBodyData {
				line := strings.Split(lines[i], ": ")
				if len(line) > 1 {
					headers[line[0]] = line[1]
				} else {
					isBodyData = true
				}
			}
			if isBodyData {
				body += lines[i]
			}

		}
		request = Request{head[0], u, head[2], headers, body}
	}
	return
}
