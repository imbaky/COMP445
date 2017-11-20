package server

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
)

//Server holds the pertinent information
type server struct {
	Host       string
	Port       string
	Protocol   string
	PathString string
	dir        path
	Debug      bool
}

//Server returns the server struct
func Server(host, port, protocol, path string, debug bool) (*server, error) {
	if host == "" {
		host = "localhost"
	}
	if port == "" {
		host = "8080"
	}
	if protocol == "" {
		protocol = "tcp"
	}
	if path == "" {
		path = "."
	}
	dir, err := os.Open(s.Dir)
	if err != nil {
		return nil, err
	}
	defer dir.Close()
	if fi, err := dir.Stat(); err != nil || fi.IsDir() {
		return nil, err
	}
	return server{host, port, protocol, dir, path, debug}, err
}

//Run start the server. It begins to listen for connections and handles those connections
func (s *server) Run() error {

	listener, err := net.Listen(s.Protocol, s.Host+":"+s.Port)
	if err != nil {
		return err
	}
	// Close the listener when the application closes.
	defer listener.Close()
	fmt.Println("listening on " + s.Host + ":" + s.Port)

	for {
		// Listen for an incoming connection.
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		// Handle connections in a new goroutine.
		go handleRequest(conn)
	}
}

// Handles incoming requests.
func handleRequest(conn net.Conn) {
	defer conn.Close()
	request, err := ioutil.ReadAll(conn)
	if err != nil {
		fmt.Printf("could not read connection %v", err)
		return
	}
}
