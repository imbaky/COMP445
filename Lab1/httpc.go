package main

import (
	"net/url"
	// "bufio"
	"io/ioutil"
	"fmt"
	"net"
	"os"
)

func main() {
	// Execution of CLI tool behavior goes here
	helpStr := `httpc is a curl like application but supports HTTP protocol only.
	Usage:	 httpc command [arguments]
	The commands are:
		get 	executes a HTTP GET request and prints the response.
		post 	executes a HTTP POST request and prints the response.
		help 	prints this screen.`

	helpGetStr := `usage: httpc get [-v] [-h key:value] URL
	Get executes a HTTP GET request for a given URL.
	  -v   		 Prints the detail of the response such as a protocol, status, and headers.
	  -h key:value   Associates headers to HTTP Request with the format 'key:value'.`

	helpPostStr := `usage: httpc post [-v] [-h key:value] [-d inline-data] [-f file] URL
	
	Post executes a HTTP POST request for a given URL with inline data or from file.
	   -v		Prints the detail of the response such as protocol, status, and headers.
	   -h key:value	Associates headers to HTTP Request with the format 'key:value'.
	   -d string	Associates inline data to the body HTTP POST request.
	   -f file	Associates the content of a file to the body HTTP POST request.
	
	Either [-d] or [-f] can be used but not both.`

	if len(os.Args) == 2 && os.Args[1] == "help" {
		fmt.Println(helpStr)
	}

	if len(os.Args) == 3 && os.Args[1] == "help" && os.Args[2] == "get" {
		fmt.Println(helpGetStr)
	}

	if len(os.Args) == 3 && os.Args[1] == "help" && os.Args[2] == "post" {
		fmt.Println(helpPostStr)
	}

	if len(os.Args) == 3 && os.Args[1] == "get" {
		u, err := url.Parse("http://"+os.Args[2])
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not parse url: %v", err)
		}
		//resolve address
		addr, err := net.ResolveTCPAddr("tcp", u.Host+":80")
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not resolve address: %v", err)
		}

		//open a connection
		conn, err := net.DialTCP("tcp", nil, addr)
		defer conn.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not open connection: %v", err)
		}

		//write to connection
		fmt.Fprintf(conn, "GET /%s HTTP/1.0\r\nHost: www.%s\r\n\r\n",u.Path,u.Host)

		//read from connection
		result, err := ioutil.ReadAll(conn)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading from connection :%v", err)
		}
		fmt.Printf("%s", result)
	}
}
