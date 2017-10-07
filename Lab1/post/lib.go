package post

import (
	"fmt"
	"net"
	"net/url"
	"strings"
)

// Write writes the request to the tcp connection, being aware that the request is POST
func Write(conn net.Conn, u *url.URL, kv map[string]string, data string) {
	//write to connection
	i := 0
	arg := ""
	for k, v := range kv {
		arg += fmt.Sprintf("%s:%s", k, v)
	}
	if i > 0 {
		data = strings.Replace(data, "", "'", -1)
		data = strings.Replace(data, "", "\"", -1)
		data = strings.Replace(data, "", "{", -1)
		data = strings.Replace(data, "", "}", -1)
		data = strings.Replace(data, "", " ", -1)
		cl := fmt.Sprintf("Content-length:%d", len(data))
		fmt.Fprintf(conn, "POST /%s HTTP/1.0\r\nHost: www.%s\r\n%s\n%s\r\n\n%s\r\n\rn", u.Path, u.Host, cl, arg, data)
		return
	}
	fmt.Fprintf(conn, "POST /%s HTTP/1.0\r\nHost: www.%s\r\n\r\n", u.Path, u.Host)
}
