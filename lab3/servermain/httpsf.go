package main

import (
	"bufio"
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/imbaky/COMP445/lab3/lab2"
	"github.com/imbaky/COMP445/lab3/server"
)

//Flag variables
var v bool
var p int
var d string

// Prints logs when in verbose mode
func verbose(log string) {
	if v {
		fmt.Println(log)
	}
}

func main() {

	// cli tool usage info message to be displayed
	helpStr := `httpfs is a simple file server.
		usage: httpfs [-v] [-p PORT] [-d PATH-TO-DIR]
		 -v Prints debugging messages.
		 -p Specifies the port number that the server will listen and serve at.
		 Default is 8080.
		 -d Specifies the directory that the server will use to read/write
		requested files. Default is the current directory when launching the
		application.`

	flag.BoolVar(&v, "v", false, "Prints debugging messages")
	flag.StringVar(&d, "d", ".", "Specifies the directory that the server will use to read/write requested files. Default is the current directory when launching the application.")
	flag.IntVar(&p, "p", 8080, "Specifies the port number that the server will listen and serve at. Default is 8080")
	flag.Parse()

	//Checks if help command is present
	var help bool
	for _, arg := range os.Args {
		if arg == "help" {
			help = true
		}
	}

	// if help flag is present
	if help {
		fmt.Println(helpStr)
	} else {

		port := fmt.Sprintf(":%d", p)

		verbose(fmt.Sprintf("Listening to Localhost%s", port))
		reqChan := make(chan server.RequestConnection)
		// listen on all interfaces
		server.Listen("localhost", port, 1000, reqChan)
		handleRequest((<-reqChan).Request, (<-reqChan).Connection)
	}

}

// Handles server connections
func handleRequest(request server.Request, conn *server.Connection) {
	fmt.Printf("Request received from %v:%v\n", conn.Peer, conn.Port)

	var response server.Response
	response = server.Response{request.Httpversion, "OK", "200", make(map[string]string), ""}
	request.URL.Path = strings.Replace(request.URL.Path, "..", "", -1)

	// GET method
	if request.Method == "GET" {
		var a []lab2.File
		response.Headers["Content-Type"] = "text/plain"
		response.Headers["Content-Disposition"] = "inline"
		isFile := false
		if request.URL.Path == "/" {
			files, err := ioutil.ReadDir(d + request.URL.Path)
			if err != nil {
				response.Error = "404"
				response.Status = "Home directory does not exist, please contact server admin"
				log.Println(err)
			}
			for _, f := range files {
				a = append(a, lab2.File{f.Name(), ""})
			}
		} else {
			efile, err := ioutil.ReadFile(d + request.URL.Path)
			fmt.Println(fmt.Sprintf("%s", efile))
			a = append(a, lab2.File{request.URL.Path, fmt.Sprintf("%s", efile)})
			isFile = true
			if err != nil {
				response.Error = "404"
				response.Status = "Requested file does not exist"
				log.Println(err)
			}
		}

		var body string

		switch request.Headers["accept"] {
		case "application/json":
			jsonData, _ := json.Marshal(a)
			body = fmt.Sprintf("%s", jsonData)
			verbose(body)
			response.Headers["Content-Type"] = "application/json"
			break
		case "text/html":
			var htmlData string
			for _, element := range a {
				htmlData += fmt.Sprintf("<li>%s</li>\n", element.FileName)
			}
			body = fmt.Sprintf("<html>\n<body>\n<ul>\n%s</ul>\n</body>\n</html>\n", htmlData)
			verbose(body)
			response.Headers["Content-Type"] = "text/html"
			break
		case "text/xml":
			xmlData, _ := xml.Marshal(a)
			body = fmt.Sprintf("%s", xmlData)
			verbose(body)
			response.Headers["Content-Type"] = "text/xml"
			break
		default:
			for _, element := range a {
				body += element.FileName + " \n" + element.Content + "\n"
			}
			verbose(body)
			break
		}
		if response.Error == "200" {
			if isFile {
				response.Headers["Content-Disposition"] = "attachment; filename=\"" + a[0].FileName + "\""
			}
			response.Body = body
		}

	}
	// POST method
	if (request.Method == "POST") && request.URL.Path != "/" {
		// write the whole body at once
		// get FileInfo structure describing file
		_, err := os.Stat(d + request.URL.Path)
		if os.IsNotExist(err) {
			fileHandle, _ := os.Create(d + request.URL.Path)
			writer := bufio.NewWriter(fileHandle)
			fmt.Fprintln(writer, request.Body)
			writer.Flush()
			fileHandle.Close()
		} else {
			if request.URL.Query()["overwrite"] != nil {
				if request.URL.Query()["overwrite"][0] == "true" {

					verbose("Overwrite is true")

					fileHandle, _ := os.Create(d + request.URL.Path)
					writer := bufio.NewWriter(fileHandle)
					fmt.Fprintln(writer, request.Body)
					writer.Flush()
					fileHandle.Close()
					response.Body = "File successfully written!!"
				} else {
					response.Status = "File exists but cannot be re-written"
					response.Error = "401"
				}
			} else {
				response.Status = "File exists but cannot be re-written"
				response.Error = "401"
			}

		}

	} else {
		if (request.Method == "POST") && request.URL.Path == "/" {
			response.Error = "400"
			response.Status = "Bad Request"
		}

	}

	verbose(fmt.Sprint(response.ToString() + "\r\n"))

	server.Respond(response, conn)
}
