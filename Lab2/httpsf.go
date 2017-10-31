package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

var v bool
var p int
var d string

func response(w http.ResponseWriter, r *http.Request) {

	if v {
		fmt.Printf("%v %v %v \n", r.Method, r.URL, r.Proto)
		fmt.Printf("Host: %v\n", r.Host)
		// Loop through headers
		for name, headers := range r.Header {
			name = strings.ToLower(name)
			for _, h := range headers {
				fmt.Printf("%v: %v \n", name, h)
			}
		}
	}
	r.ParseForm() // parse arguments

	if r.Method == "GET" {
		files, err := ioutil.ReadDir("." + r.URL.Path)
		if err != nil {
			http.Error(w, "directory not found", http.StatusNotFound)
			log.Println(err)
		}

		for _, f := range files {
			fmt.Fprintf(w, f.Name()+"\n")
		}
	}
	if (r.Method == "POST") && r.URL.Path != "/" {
		// write the whole body at once
		bodyBytes, _ := ioutil.ReadAll(r.Body)
		bodyString := string(bodyBytes)

		fileHandle, _ := os.Create("./" + r.URL.Path + ".txt")
		writer := bufio.NewWriter(fileHandle)

		fmt.Fprintln(writer, bodyString)
		writer.Flush()
		fileHandle.Close()
		fmt.Fprintf(w, bodyString)

	}

}

func main() {
	helpStr := `httpfs is a simple file server.
	usage: httpfs [-v] [-p PORT] [-d PATH-TO-DIR]
	 -v Prints debugging messages.
	 -p Specifies the port number that the server will listen and serve at.
	 Default is 8080.
	 -d Specifies the directory that the server will use to read/write
	requested files. Default is the current directory when launching the
	application.`

	flag.BoolVar(&v, "v", false, "Prints debugging messages")
	flag.StringVar(&d, "d", "./", "Specifies the directory that the server will use to read/write requested files. Default is the current directory when launching the application.")
	flag.IntVar(&p, "p", 8080, "Specifies the port number that the server will listen and serve at. Default is 8080")
	flag.Parse()

	argsmap := make(map[string]bool)
	for _, arg := range os.Args {
		argsmap[arg] = true
	}
	if argsmap["help"] {
		fmt.Println(helpStr)
	} else {
		http.HandleFunc("/", response) // set router
		port := fmt.Sprintf(":%d", p)
		if v {
			fmt.Printf("Listening to Localhost%s\n", port)
		}
		err := http.ListenAndServe(port, nil) // set listen port
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}

	}

}
