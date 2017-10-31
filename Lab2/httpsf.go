package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func response(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // parse arguments
	fmt.Println()
	if r.Method == "GET" {
		fmt.Println("path", r.URL.Path)

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

	var v bool
	var p int
	var d string

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
		http.HandleFunc("/", response)           // set router
		err := http.ListenAndServe(":9090", nil) // set listen port
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}

}
