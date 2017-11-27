package server

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
	"net/url"
	"strings"
	"time"

	"github.com/imbaky/COMP445/lab3/packet"
)

// File struct definition
type File struct {
	FileName string
	Content  string
}

// Request struct definition
type Request struct {
	Method      string
	URL         *url.URL
	Httpversion string
	Headers     map[string]string
	Body        string
}

// Response struct definition
type Response struct {
	HTTPVersion string
	Status      string
	Error       string
	Headers     map[string]string
	Body        string
}

type Connection struct {
	Conn     *net.UDPConn
	Timeout  int
	Sequence uint32
	WindowK  int
	Buffer   []packet.Packet
}

func Listen(host, port string, timeout, windowK int) (*Connection, error) {
	addr, err := net.ResolveUDPAddr("udp", port)
	if err != nil {
		return nil, err
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, err
	}
	bufSize := 1
	for i := 1; 0 < windowK; i++ {
		bufSize *= 2
	}
	buffer := make([]packet.Packet, bufSize)
	connection := Connection{conn, timeout, 0, windowK, buffer}
	return &connection, nil
}

func (conn *Connection) readPacket() packet.Packet {
	buf := make([]byte, 1024)
	n, _, _ := conn.Conn.ReadFromUDP(buf)
	pkt, _ := packet.FromBytes(buf[0:n])
	return pkt
}

func (conn *Connection) Write(pkt packet.Packet) error {
	_, err := conn.Conn.Write(pkt.Bytes())
	return err
}
func (conn *Connection) generateSequence() {
	conn.Sequence = uint32(rand.Intn(len(conn.Buffer)))
}

func getUint32(buff []byte) uint32 {
	return binary.LittleEndian.Uint32(buff)
}

func Establish(conn *Connection) bool {
	// make the packet channel to read and write to for the handshake
	c := make(chan packet.Packet, 1)
	pkt := conn.readPacket()
	if pkt.PType != packet.SYN { // did not get an establish SYN packet
		return false
	}
	clientSize := binary.LittleEndian.Uint32(pkt.Payld) // this is the clients windowk
	conn.setWindowAndBuffer(int(clientSize))
	conn.generateSequence()
	windowk := make([]byte, 4)
	seq := make([]byte, 4)
	binary.LittleEndian.PutUint32(seq, conn.Sequence)
	binary.LittleEndian.PutUint32(windowk, uint32(conn.WindowK))
	synack, err := packet.MakePacket(packet.SYNACK, seq, pkt.Peer, pkt.Port, windowk) // Send back the nack with the seq number and the final windowk
	if err != nil {
		return false
	}
	err = conn.Write(synack)
	if err != nil {
		return false
	}
	for { // keep reading and if not ack send back synack
		go func(conn *Connection) {
			c <- conn.readPacket()
		}(conn)

		select {
		case <-time.After(time.Millisecond * time.Duration(conn.Timeout)):
			conn.Write(synack)
		case res := <-c:
			if res.PType == packet.ACK && getUint32(res.Seq) == (conn.Sequence+1) {
				return true
			}
			conn.Write(synack)
		}
	}
	return false
}

func (conn *Connection) setWindowAndBuffer(windowK int) {
	if conn.WindowK > windowK {
		conn.WindowK = windowK
	}
	bufSize := 1
	for i := 1; 0 < conn.WindowK; i++ {
		bufSize *= 2
	}
	conn.Buffer = make([]packet.Packet, bufSize)
}

//converts response to string
func (response Response) ToString() (responseText string) {
	responseText = fmt.Sprintf("%s %s %s \r\n", response.HTTPVersion, response.Error, response.Status)
	response.Headers["Server"] = "COMP445/2.0 (Assignment)"
	now := time.Now()
	response.Headers["Last-Modified"] = now.Format("Mon Jan 2 15:04:05 MST 2006")
	response.Headers["Content-Length"] = fmt.Sprintf("%d", len(response.Body))
	for name, value := range response.Headers {
		responseText += fmt.Sprintf("%s: %s \r\n", name, value)
	}
	responseText += fmt.Sprintf("%s \r\n\r\n", response.Body)
	return
}

func formRequest(conn *net.UDPConn) []byte {

}

//Function to take buffer data and parse it into Request
func ParseRequest(buf []byte) (request Request) {
	var lines []string
	//Split each line
	lines = strings.Split(string(buf), "\r\n")

	// Split the first line with the request definition
	head := strings.Split(lines[0], " ")

	if len(head) > 2 {
		u, _ := url.Parse(head[1])

		headers := make(map[string]string)
		var body string
		isBodyData := false
		for i := 1; i < len(lines); i += 2 {
			if lines[i] != "" && !isBodyData {
				line := strings.Split(lines[i], ": ")
				if len(line) > 1 {
					headers[line[0]] = line[1]
				} else {
					isBodyData = true
				}
			}
			if isBodyData {
				body += lines[i]
			}

		}
		request = Request{head[0], u, head[2], headers, body}
	}
	return
}
