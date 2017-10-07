package get

import (
	"fmt"
	"net"
	"net/url"
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
	fmt.Println(u.Path)
	fmt.Fprintf(conn, "GET /%s HTTP/1.0\r\nHost: www.%s\r\n\r\n", u.Path, u.Host)
}
