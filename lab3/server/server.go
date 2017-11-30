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

const (
	buffSize = 100000
)

// File struct definition
type File struct {
	FileName string
	Content  string
}
type RequestConnection struct {
	Request    Request
	Connection *Connection
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
	Buffer   []*packet.Packet
	Peer     []byte
	Port     []byte
}

func Respond(response Response, conn *Connection) {
	respStr := response.ToString()
	conn.Buffer = []*packet.Packet{}
	i := 0
	j := 0
	for ; i < len([]byte(respStr)); i += 1014 {
		seq := make([]byte, 4)
		binary.BigEndian.PutUint32(seq, uint32(j))
		payload := []byte{}
		copy([]byte(respStr)[i:i+1014], payload)
		pkt, _ := packet.MakePacket(packet.ACK, seq, conn.Peer, conn.Port, payload)
		conn.Buffer[j] = &pkt
		j++
	}
	seq := make([]byte, 4)
	binary.BigEndian.PutUint32(seq, uint32(j))
	payload := []byte{}
	copy([]byte(respStr)[i:], payload)
	pkt, _ := packet.MakePacket(packet.ACK, seq, conn.Peer, conn.Port, payload)
	conn.Buffer[j] = &pkt
	establishResponse(conn)
	SendResponse(conn)
}

func SendResponse(conn *Connection) {
	acks := make([]bool, len(conn.Buffer))
	chacks := make(chan []bool, 1)
	total := len(conn.Buffer)
	chacks <- acks
	for {
		select {
		case <-time.After(time.Millisecond * time.Duration(conn.Timeout)):
			if total == 0 {
				return
			}
			ackArr := <-chacks
			for k, v := range acks {
				if !v {
					conn.Write(*conn.Buffer[k])
				}
			}
			chacks <- ackArr
		case ackArr := <-chacks:
			pkt := conn.readPacket()
			seq := getUint32(pkt.Seq)
			if !ackArr[int(seq)] {
				total--
				ackArr[int(seq)] = true
			}
			conn.Write(*(conn.Buffer[int(seq)]))
			chacks <- ackArr
		}
	}

}

func Listen(host, port string, timeout int, ch chan<- RequestConnection) error {
	addr, err := net.ResolveUDPAddr("udp", port)
	if err != nil {
		return err
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return err
	}
	buffer := make([]*packet.Packet, buffSize)
	connection := Connection{conn, timeout, 0, buffer, []byte{}, []byte{}}

	for {
		establish(&connection)
		ch <- RequestConnection{ParseRequest(receive(&connection)), &connection}
	}
}

func establishResponse(conn *Connection) {
	syn, _ := packet.MakePacket(packet.SYN, []byte{}, conn.Peer, conn.Port, []byte{}) // Send back the nack with the seq number and the final windowk
	binary.BigEndian.PutUint32(syn.Seq, uint32(len(conn.Buffer)))
	conn.Write(syn)
	conn.Write(syn)
	conn.Write(syn)
}

func checkTimeout(conn *Connection) bool {
	var peer, port, seq []byte
	done := true
	for _, v := range conn.Buffer {
		if v != nil {
			peer, port, seq = v.Peer, v.Port, v.Seq
			break
		}
	}
	for _, v := range conn.Buffer {
		if v == nil {
			done = false
			nack, _ := packet.MakePacket(packet.NACK, seq, peer, port, []byte{})
			conn.Write(nack)
		}
	}
	return done
}

func receive(conn *Connection) []byte {
	total := len(conn.Buffer)
	go func() {
		for {
			time.Sleep(time.Millisecond * time.Duration(conn.Timeout))
			if checkTimeout(conn) {
				break
			}
		}
	}()

	for total > 0 {
		pkt := conn.readPacket()
		if conn.Buffer[pkt.GetSequence()] != nil {
			ack, _ := packet.MakePacket(packet.ACK, pkt.Seq, pkt.Peer, pkt.Port, []byte{})
			conn.Write(ack)
		} else {
			copy := pkt
			conn.Buffer[pkt.GetSequence()] = &copy
			ack, _ := packet.MakePacket(packet.ACK, pkt.Seq, pkt.Peer, pkt.Port, []byte{})
			conn.Write(ack)
		}
	}
	return extractPayloads(conn.Buffer)
}

func extractPayloads(pkts []*packet.Packet) []byte {
	buf := []byte{}
	for _, v := range pkts {
		buf = append(buf, v.Payld...)
	}
	return buf
}

func establish(conn *Connection) {
	for {
		pkt := conn.readPacket()
		if pkt.PType == packet.SYN { // did not get an establish SYN packet
			conn.Buffer = make([]*packet.Packet, pkt.GetSequence())
			seq := []byte{0x00, 0x00, 0x00, 0x00}
			synack, _ := packet.MakePacket(packet.SYNACK, seq, pkt.Peer, pkt.Port, []byte{}) // Send back the nack with the seq number and the final windowk
			conn.Write(synack)
			conn.Write(synack)
			conn.Write(synack)
			return
		}
	}
	return
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
