package main

import (
	"bufio"
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
)

// File struct definition
type File struct {
	FileName string
	content  string
}

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
	flag.StringVar(&d, "d", "./", "Specifies the directory that the server will use to read/write requested files. Default is the current directory when launching the application.")
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

	buf := make([]byte, 1024)
	n, re := conn.Read(buf)
	if re != nil {
		fmt.Fprintf(os.Stderr, "read error %v\n", re)
	}

	verbose(fmt.Sprintln("\n", string(buf)))

	// Parsing request into Request struct
	var request Request
	if n > 0 {
		request = parseRequest(buf)
	} else {
		verbose("Invalid Request")
		return
	}

	var response Response
	response = Response{request.httpversion, "OK", "200", make(map[string]string), ""}

	// GET method
	if request.method == "GET" {
		var a []File
		if request.URL.Path == "/" {
			files, err := ioutil.ReadDir(d + request.URL.Path)
			if err != nil {
				response.Error = "404"
				log.Println(err)
			}
			for _, f := range files {
				a = append(a, File{f.Name(), ""})
			}
		} else {
			efile, err := ioutil.ReadFile(d + request.URL.Path)
			fmt.Println(fmt.Sprintf("%s", efile))
			a = append(a, File{"." + request.URL.Path, fmt.Sprintf("%s", efile)})
			if err != nil {
				response.Error = "404"
				log.Println(err)
			}
		}

		var body string
		switch request.headers["accept"] {
		case "application/json":
			jsonData, _ := json.Marshal(a)
			body = fmt.Sprintf("%s", jsonData)
			verbose(body)
			break
		case "text/html":
			var htmlData string
			for _, element := range a {
				htmlData += fmt.Sprintf("<li>%s</li>\n", element.FileName)
			}
			body = fmt.Sprintf("<html>\n<body>\n<ul>\n%s</ul>\n</body>\n</html>\n", htmlData)
			verbose(body)
			break
		case "text/xml":
			xmlData, _ := xml.Marshal(a)
			body = fmt.Sprintf("%s", xmlData)
			verbose(body)
			break
		default:
			for _, element := range a {
				body += element.FileName + " \n" + element.content + "\n"
			}
			verbose(body)
			break
		}
		response.Body = body

	}
	// POST method
	if (request.method == "POST") && request.URL.Path != "/" {
		// write the whole body at once
		// get FileInfo structure describing file
		_, err := os.Stat(d + request.URL.Path + ".txt")
		if os.IsNotExist(err) {
			fileHandle, _ := os.Create(d + request.URL.Path + ".txt")
			writer := bufio.NewWriter(fileHandle)
			fmt.Fprintln(writer, request.body)
			writer.Flush()
			fileHandle.Close()
		} else {
			if request.URL.Query()["overwrite"] != nil {
				if request.URL.Query()["overwrite"][0] == "true" {

					verbose("Overwrite is true")

					fileHandle, _ := os.Create(d + request.URL.Path + ".txt")
					writer := bufio.NewWriter(fileHandle)
					fmt.Fprintln(writer, request.body)
					writer.Flush()
					fileHandle.Close()
				}
			}

		}
		response.Body = request.body

	}

	verbose(fmt.Sprint(response.toString() + "\r\n"))

	if _, we := conn.Write([]byte(response.toString())); we != nil {
		fmt.Fprintf(os.Stderr, "write error %v\n", we)
	}
	conn.Close()

}
