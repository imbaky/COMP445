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
	"net/url"
	"os"
	"strings"
)

/*
*	File struct definition
 */
type File struct {
	FileName string
	content  string
}

/*
*	Request struct definition
 */
type Request struct {
	method      string
	URL         *url.URL
	httpversion string
	headers     map[string]string
	body        string
}

//Flag variables
var v bool
var p int
var d string

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
		port := fmt.Sprintf(":%d", p)

		if v {
			fmt.Printf("Listening to Localhost%s\n", port)
		}
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

func handleConn(conn net.Conn) {
	defer conn.Close()
	fmt.Printf("Request recieved from %v\n", conn.RemoteAddr())

	buf := make([]byte, 1024)
	n, re := conn.Read(buf)
	if re != nil {
		fmt.Fprintf(os.Stderr, "read error %v\n", re)
	}

	if v {
		fmt.Println("\n", string(buf))
	}

	var request Request
	if n > 0 {
		var lines []string
		lines = strings.Split(string(buf), "\r\n")
		head := strings.Split(lines[0], " ")

		if len(head) > 2 {
			u, _ := url.Parse(head[1])
			headers := make(map[string]string)
			var body string
			for i := 1; i < len(lines); i += 2 {
				if lines[i] != "" {
					line := strings.Split(lines[i], ":")
					headers[line[0]] = line[1]
				} else {
					if i < len(lines) {
						body = lines[i+1]
					}
					break
				}
				request = Request{head[0], u, head[2], headers, body}
			}

		}

		var response string

		if request.method == "GET" {
			var a []File
			if request.URL.Path == "/" {
				files, err := ioutil.ReadDir("." + request.URL.Path)
				if err != nil {
					response += request.httpversion + " 404 OK \r\n"
					log.Println(err)
				}
				response += request.httpversion + " 200 OK \r\n"
				for _, f := range files {
					a = append(a, File{f.Name(), ""})
				}
			} else {
				response += request.httpversion + " 200 OK \r\n"
				efile, err := ioutil.ReadFile("." + request.URL.Path)
				fmt.Println(fmt.Sprintf("%s", efile))
				a = append(a, File{"." + request.URL.Path, fmt.Sprintf("%s", efile)})
				if err != nil {
					response += request.httpversion + " 404 OK \r\n"
					log.Println(err)
				}
			}

			var body string
			switch request.headers["accept"] {
			case "application/json":
				jsonData, _ := json.Marshal(a)
				body = fmt.Sprintf("%s", jsonData)
				fmt.Println(body)
				break
			case "text/html":
				var htmlData string
				for _, element := range a {
					htmlData += fmt.Sprintf("<li>%s</li>\n", element.FileName)
				}
				body = fmt.Sprintf("<html>\n<body>\n<ul>\n%s</ul>\n</body>\n</html>\n", htmlData)
				break
			case "text/xml":
				xmlData, _ := xml.Marshal(a)
				body = fmt.Sprintf("%s", xmlData)
				break
			default:
				for _, element := range a {
					body = element.FileName + " \n" + element.content + "\n"
				}
				break
			}
			response += "Message Body:\n"
			response += body + "\r\n"

		}
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
						if v {
							fmt.Println("Overwrite is true")
						}
						fileHandle, _ := os.Create(d + request.URL.Path + ".txt")
						writer := bufio.NewWriter(fileHandle)
						fmt.Fprintln(writer, request.body)
						writer.Flush()
						fileHandle.Close()
					}
				}

			}

			response += request.body

		}
		if v {
			fmt.Println(response + "\r\n")
		}
		if _, we := conn.Write([]byte(response)); we != nil {
			fmt.Fprintf(os.Stderr, "write error %v\n", we)
		}
		conn.Close()
	}

}
