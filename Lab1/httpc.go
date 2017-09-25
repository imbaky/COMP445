package main

import (
	"fmt"
	"net"
	"os"
	"bufio"
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
		conn, err := net.Dial("tcp", os.Args[2])
		if err != nil {
			fmt.Printf("error: %v\n",err)
			return
		}
		fmt.Fprintf(conn, "GET / HTTP/1.0\r\nHost: "+os.Args[2]+"\r\n\r\n")
		status, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Printf("error: %v\n",err)
			return
		}
		fmt.Printf("status: %v\n",status)
	}
}
