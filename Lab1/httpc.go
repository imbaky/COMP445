package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/url"
	"os"
	"strings"

	"github.com/imbaky/COMP445/Lab1/get"
	"github.com/imbaky/COMP445/Lab1/post"
)

func main() {
	// Execution of CLI tool behavior goes here
	helpStr := `httpc is a curl like application but supports HTTP protocol only.
	Usage:	 httpc command [arguments]
	The commands are:
		get 	executes a HTTP GET request and prints the response.
		post 	executes a HTTP POST request and prints the response.
		help 	prints this screen.`

	helpGetStr := `usage: httpc get [-v] [-h key1:value1,key2,value2] URL
	Get executes a HTTP GET request for a given URL.
	  -v   		 Prints the detail of the response such as a protocol, status, and headers.
	  -h key:value   Associates headers to HTTP Request with the format 'key:value'.`

	helpPostStr := `usage: httpc post [-v] [-h key1:value1,key2,value2] [-d inline-data] [-f file] URL
	
	Post executes a HTTP POST request for a given URL with inline data or from file.
	   -v		Prints the detail of the response such as protocol, status, and headers.
	   -h key:value	Associates headers to HTTP Request with the format 'key:value'.
	   -d string	Associates inline data to the body HTTP POST request.
	   -f file	Associates the content of a file to the body HTTP POST request.
	
	Either [-d] or [-f] can be used but not both.`

	var flagList string
	var data string
	flag.StringVar(&flagList, "h", "", "key value in the format key1:value1,key2,value2")
	flag.StringVar(&data, "d", "", "key value in the format key1:value1,key2,value2")
	flag.Parse()

	kvmap := make(map[string]string)
	for _, v := range strings.Split(flagList, ",") {
		fmt.Printf("v %v\n", v)
		pair := strings.Split(v, ":")
		fmt.Printf("pair %v\n", pair)
		kvmap[pair[0]] = pair[1]
	}

	argsmap := make(map[string]bool)
	for _, arg := range os.Args {
		argsmap[arg] = true
	}

	if argsmap["help"] {
		fmt.Println(helpStr)
		if argsmap["get"] {
			fmt.Println(helpGetStr)
			return
		}
		if argsmap["post"] {
			fmt.Println(helpPostStr)
			return
		}
	}

	if argsmap["get"] {
		var index int
		for i, v := range os.Args {
			if v == "get" {
				index = i + 1
				continue
			}
		}
		u, err := url.Parse("http://" + os.Args[index])
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

		get.Write(conn, u, kvmap)
		//read from connection
		result, err := ioutil.ReadAll(conn)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading from connection :%v", err)
		}
		fmt.Printf("%s", result)
		return
	}

	if argsmap["post"] {
		var index int
		for i, v := range os.Args {
			if v == "post" {
				index = i + 1
				continue
			}
		}
		u, err := url.Parse("http://" + os.Args[index])
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

		post.Write(conn, u, kvmap,data)
		//read from connection
		result, err := ioutil.ReadAll(conn)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading from connection :%v", err)
		}
		fmt.Printf("%s", result)
	}

}
