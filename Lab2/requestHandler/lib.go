package requestHandler

import (
	"fmt"
	"net/url"
	"strings"
)

//Function to take buffer data and parse it into Request
func parseRequest(buf []byte) (request Request) {
	var lines []string
	//Split each line
	lines = strings.Split(string(buf), "\r\n")

	// Split the first line with the request definition
	head := strings.Split(lines[0], " ")

	if len(head) > 2 {
		u, _ := url.Parse(head[1])
		headers := make(map[string]string)
		var body string
		for i := 1; i < len(lines); i += 2 {
			if lines[i] != "" {
				line := strings.Split(lines[i], ": ")
				if len(line) > 1 {
					headers[line[0]] = line[1]
				}

			} else {
				body += lines[i]
			}

		}
		request = Request{head[0], u, head[2], headers, body}
	}
	return
}

// Request struct definition
type Request struct {
	method      string
	URL         *url.URL
	httpversion string
	headers     map[string]string
	body        string
}

// Response struct definition
type Response struct {
	HTTPVersion string
	Status      string
	Error       string
	Headers     map[string]string
	Body        string
}

func (response Response) toString() (responseText string) {
	responseText = fmt.Sprintf("%s %s %s \r\n", response.HTTPVersion, response.Error, response.Status)
	for name, value := range response.Headers {
		responseText += fmt.Sprintf("%s: %s \r\n", name, value)
	}
	responseText += fmt.Sprintf("%s \r\n", response.Body)
	return
}
