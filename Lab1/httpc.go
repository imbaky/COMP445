package main

import (
	"fmt"
	"os"
)

func main() {
	// Execution of CLI tool behavior goes here
	helpStr := "httpc is a curl like application but supports HTTP protocol only.\nUsage:\n\t httpc command [arguments]\n The commands are:\n\t get \t executes a HTTP GET request and prints the response.\n\t post \t executes a HTTP POST request and prints the response.\n\t help \t prints this screen.\n"
	helpGetStr := "usage: httpc get [-v] [-h key:value] URL\nGet executes a HTTP GET request for a given URL.\n  -v   \t\t Prints the detail of the response such as a protocol, status, and headers.\n  -h key:value   Associates headers to HTTP Request with the format 'key:value'.\n"
	helpPostStr := "usage: httpc post [-v] [-h key:value] [-d inline-data] [-f file] URL\n\nPost executes a HTTP POST request for a given URL with inline data or from file.\n\n\t-v\t\tPrints the detail of the response such as protocol, status, and headers.\n\t-h key:value\tAssociates headers to HTTP Request with the format 'key:value'.\n\t-d string\tAssociates inline data to the body HTTP POST request.\n\t-f file\t\tAssociates the content of a file to the body HTTP POST request.\n\nEither [-d] or [-f] can be used but not both.\n"
	if len(os.Args) == 2 && os.Args[1] == "help" {
		fmt.Print(helpStr)
	}

	if len(os.Args) == 3 && os.Args[1] == "help" && os.Args[2] == "get" {
		fmt.Print(helpGetStr)
	}

	if len(os.Args) == 3 && os.Args[1] == "help" && os.Args[2] == "post" {
		fmt.Print(helpPostStr)
	}
}
