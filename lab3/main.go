package main

import (
	"fmt"
	"net"
	"os"

	"github.com/imbaky/COMP445/lab3/packet"
)

func main() {
	// pt := packet.ACK
	// seq := []byte{0x00, 0x01, 0x02, 0x03}
	// b := make([]byte, 4)
	// littleEnd := binary.LittleEndian.Uint32(seq)
	// binary.BigEndian.PutUint32(b, littleEnd)
	// fmt.Printf("seq %x \n", seq)

	// peer := []byte{0x00, 0x01, 0x02, 0x03}
	// port := []byte{0x00, 0x01}
	// payld := []byte{}
	// pkt, err := packet.Packet(pt, seq, peer, port, payld)
	// if err != nil {
	// 	fmt.Printf("error %v", err)
	// 	return
	// }
	// fmt.Printf("packet %v\n", pkt)
	// fmt.Printf("bytes %v\n", pkt.Bytes())
	// fmt.Printf("seq %x \n", pkt.Seq)
	// fmt.Printf("seq %x \n", seq)
	// fd, _ := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, syscall.PROT_NONE)
	// f := os.NewFile(uintptr(fd), fmt.Sprintf("fd %d", fd))

	// for {
	// 	buf := make([]byte, 1024)
	// 	numRead, err := f.Read(buf)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}
	// 	fmt.Printf("% X\n", buf[:numRead])
	// }
	// pkt, host, port, err := server.ReceivePackets("127.0.0.1", ":8007")
	// if err != nil {
	// 	fmt.Printf("read packets returned %v", err)
	// }
	// fmt.Printf("packet %v\n host %v \n port %v \n", pkt, host, port)
	ServerAddr, err := net.ResolveUDPAddr("udp", ":8007")
	CheckError(err)

	/* Now listen at selected port */
	ServerConn, err := net.ListenUDP("udp", ServerAddr)
	CheckError(err)
	defer ServerConn.Close()

	buf := make([]byte, 1024)

	for {
		n, addr, err := ServerConn.ReadFromUDP(buf)
		fmt.Println("Received ", string(buf[0:n]), " from ", addr)
		pkt, err := packet.FromBytes(buf[0:n])
		CheckError(err)
		fmt.Printf("packet : %v\n", pkt)
	}
}
func CheckError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(0)
	}
}
