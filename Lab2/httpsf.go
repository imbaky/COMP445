package main

import (
	"bufio"
	"encoding/json"
	"encoding/xml"
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

type file struct {
	FileName string
}

func response(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // parse arguments
	bodyBytes, _ := ioutil.ReadAll(r.Body)
	bodyString := string(bodyBytes)
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
		fmt.Println("Body: " + bodyString)
	}

	if r.Method == "GET" {
		files, err := ioutil.ReadDir("." + r.URL.Path)
		if err != nil {
			http.Error(w, "directory not found", http.StatusNotFound)
			log.Println(err)
		}
		var a []file

		for _, f := range files {
			a = append(a, file{f.Name()})
		}
		switch r.Header.Get("accept") {
		case "application/json":
			jsonData, _ := json.Marshal(a)
			fmt.Fprintf(w, fmt.Sprintf("%s", jsonData))
			break
		case "text/html":
			var htmlData string
			for _, element := range a {
				htmlData += fmt.Sprintf("<li>%s</li>\n", element.FileName)
			}
			fmt.Fprintf(w, fmt.Sprintf("<html>\n<body>\n<ul>\n%s</ul>\n</body>\n</html>\n", htmlData))
			break
		case "text/xml":
			xmlData, _ := xml.Marshal(a)
			fmt.Fprintf(w, fmt.Sprintf("%s", xmlData))
			break
		default:
			for _, element := range a {
				fmt.Fprintf(w, element.FileName+"\n")
			}
			break
		}

	}
	if (r.Method == "POST") && r.URL.Path != "/" {
		// write the whole body at once
		// get FileInfo structure describing file
		_, err := os.Stat(d + r.URL.Path + ".txt")
		if os.IsNotExist(err) {
			fileHandle, _ := os.Create(d + r.URL.Path + ".txt")
			writer := bufio.NewWriter(fileHandle)
			fmt.Fprintln(writer, bodyString)
			writer.Flush()
			fileHandle.Close()
		} else {
			if r.FormValue("overwrite") == "true" {
				fmt.Println("true")
				fileHandle, _ := os.Create(d + r.URL.Path + ".txt")
				writer := bufio.NewWriter(fileHandle)
				fmt.Fprintln(writer, bodyString)
				writer.Flush()
				fileHandle.Close()
			}

		}

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
