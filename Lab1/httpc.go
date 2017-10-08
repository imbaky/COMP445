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

	var h string
	var d string
	var o string
	var v bool
	flag.StringVar(&h, "h", "", "key value in the format key1:value1,key2,value2")
	flag.BoolVar(&v, "v", false, "output the return status only")
	flag.StringVar(&d, "d", "", "inline-data")
	flag.StringVar(&o, "o", "", "file to write to")
	flag.Parse()

	kvmap := make(map[string]string)
	if h != "" {
		for _, v := range strings.Split(h, ",") {
			fmt.Printf("v %v\n", v)
			pair := strings.Split(v, ":")
			fmt.Printf("pair %v\n", pair)
			kvmap[pair[0]] = pair[1]
		}
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

	if argsmap["get"] == true {

		u, err := url.Parse("http://" + os.Args[len(os.Args)-1])
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
		status, result, err := get.Read(conn)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading from connection :%v", err)
		}
		if findArg(os.Args, "-o") {
			o = os.Args[findArgPos(os.Args, "-o")+1]
			file, err := os.OpenFile(o, os.O_RDWR|os.O_CREATE, 0755)
			defer file.Close()
			if err != nil {
				fmt.Fprintf(os.Stderr, "could not open file %s :%v", o, err)
				return
			}
			fmt.Fprint(file, result)
			return
		}
		if v || findArg(os.Args, "-v") {
			fmt.Printf("Output:\n\n%s \r\n\n%s", status, result)
		} else {
			fmt.Printf("Output:\n\n%s", result)
		}
		return
	}

	if argsmap["post"] {

		u, err := url.Parse("http://" + os.Args[len(os.Args)-1])
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

		post.Write(conn, u, kvmap, d)
		//read from connection
		result, err := ioutil.ReadAll(conn)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading from connection :%v", err)
		}
		fmt.Printf("%s", result)
	}

}

func findArg(args []string, s string) bool {
	for _, v := range args {
		if v == s {
			return true
		}
	}
	return false
}

func findArgPos(args []string, s string) int {
	for k, v := range args {
		if v == s {
			return k
		}
	}
	return 0
}
