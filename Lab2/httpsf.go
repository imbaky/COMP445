package main

import (
	"./http"
	"bufio"
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
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

		// listen on all interfaces
		listener, _ := net.Listen("tcp", port)
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Fprintf(os.Stderr, "error occured during accept connection %v\n", err)
				continue
			}
			go handleConn(conn)
		}
	}

}

// Handles server connections
func handleConn(conn net.Conn) {
	defer conn.Close()
	fmt.Printf("Request recieved from %v\n", conn.RemoteAddr())

	buf := make([]byte, 1028)
	n, re := conn.Read(buf)
	if re != nil {
		fmt.Fprintf(os.Stderr, "read error %v\n", re)
	}

	verbose(fmt.Sprintln("\n", string(buf)))

	// Parsing request into Request struct
	var request http.Request
	if n > 0 {
		request = http.ParseRequest(buf)
	} else {
		verbose("Invalid Request")
		return
	}
	//Create a response
	var response http.Response
	//intializing the response
	response = http.Response{request.Httpversion, "OK", "200", make(map[string]string), ""}

	//remove parent directory ".." in path to prevent user from accessing anything not in this directory
	request.URL.Path = strings.Replace(request.URL.Path, "..", "", -1)

	// GET method
	if request.Method == "GET" {
		var a []http.File
		response.Headers["Content-Type"] = "text/plain"
		response.Headers["Content-Disposition"] = "inline"
		isFile := false
		if request.URL.Path == "/" {
			files, err := ioutil.ReadDir(d + request.URL.Path)
			if err != nil {
				response.Error = "404"
				response.Body = "Home directory does not exist, please contact server admin"
				log.Println(err)
			}
			for _, f := range files {
				a = append(a, http.File{f.Name(), ""})
			}
		} else {
			efile, err := ioutil.ReadFile(d + request.URL.Path)
			fmt.Println(fmt.Sprintf("%s", efile))
			a = append(a, http.File{request.URL.Path, fmt.Sprintf("%s", efile)})
			isFile = true
			if err != nil {
				response.Error = "404"
				response.Body = "Requested file does not exist"
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
					response.Body = request.Body
				} else {
					response.Body = "File exists but cannot be re-written"
					response.Error = "401"
				}
			} else {
				response.Body = "File exists but cannot be re-written"
				response.Error = "401"
			}

		}

	} else {
		if (request.Method == "POST") && request.URL.Path == "/" {
			response.Error = "400"
			response.Body = "Bad Request"
		}

	}

	verbose(fmt.Sprint(response.ToString() + "\r\n"))

	if _, we := conn.Write([]byte(response.ToString())); we != nil {
		fmt.Fprintf(os.Stderr, "write error %v\n", we)
	}
	conn.Close()

}
