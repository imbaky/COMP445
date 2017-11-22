package main

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/imbaky/COMP445/lab3/packet"
)

func main() {
	pt := packet.ACK
	seq := []byte{0x00, 0x01, 0x02, 0x03}
	test := [4]byte{0x00, 0x01, 0x02, 0x03}
	tt := [4]byte{}
	fmt.Printf("test %x \n", test)

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, &test)
	binary.Read(buf, binary.BigEndian, &tt)

	fmt.Printf("tt %x \n", tt)
	peer := []byte{0x00, 0x01, 0x02, 0x03}
	port := []byte{0x00, 0x01}
	payld := []byte{}
	pkt, err := packet.Packet(pt, seq, peer, port, payld)
	if err != nil {
		fmt.Printf("error %v", err)
		return
	}
	fmt.Printf("packet %v\n", pkt)
	fmt.Printf("bytes %v\n", pkt.Bytes())
	fmt.Printf("seq %x \n", pkt.Sequence)
	fmt.Printf("seq %x \n", seq)
}
