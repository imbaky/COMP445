package get

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/url"
	"strings"
)

// Write writes the request to the tcp connection, being aware that the request is GET
func Write(conn net.Conn, u *url.URL, kv map[string]string) {
	//write to connection
	if len(kv) != 0 {
		u.Path += "?"
	}
	i := 0
	for k, v := range kv {
		if i == 0 {
			arg := fmt.Sprintf("%s=%s", k, v)
			u.Path += arg
			i++
		} else {
			arg := fmt.Sprintf("&%s=%s", k, v)
			u.Path += arg
		}
	}
	if u.RawQuery != "" {
		fmt.Fprintf(conn, "GET /%s HTTP/1.0\r\nHost: www.%s\r\n\r\n", u.Path+"?"+u.RawQuery, u.Host)
		return
	}
	fmt.Fprintf(conn, "GET /%s HTTP/1.0\r\nHost: www.%s\r\n\r\n", u.Path, u.Host)
}

// Read reads the response and returns the status and full output
func Read(conn net.Conn) (string, string, error) {
	//read from connection
	result, err := ioutil.ReadAll(conn)
	// fmt.Printf("result %v\n",string(result))
	if err != nil {
		return "", "", err
	}
	status := strings.Split(string(result), "\r\n\r\n")
	return status[0], status[1], nil
}
