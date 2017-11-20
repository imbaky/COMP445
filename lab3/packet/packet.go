package packet

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const (
	ACK    byte = 0xFF
	NACK   byte = 0x00
	SYN    byte = 0xF0
	SYNACK byte = 0x0F
)

type packet struct {
	pType byte
	seq   []byte // 4bytes
	peer  []byte // 4bytes
	port  []byte // 2bytes
	payld []byte // max 1014 bytes
}

func Packet(pType byte, seq, peer, port, payld []byte) (packet, error) {
	if pType != ACK && pType != NACK && pType != SYN && pType != SYNACK {
		return packet{}, fmt.Errorf("packet type must be one of the following %v %v %v %v", ACK, NACK, SYN, SYNACK)
	}
	toBig(&pType)
	if len(seq) != 4 {
		return packet{}, fmt.Errorf("seq is not 4 bytes")
	}
	toBigEnd(&seq)
	if len(peer) != 4 {
		return packet{}, fmt.Errorf("peer is not 4 bytes")
	}
	toBigEnd(&peer)
	if len(port) != 2 {
		return packet{}, fmt.Errorf("port is not 2 bytes")
	}
	toBigEnd(&port)
	if len(payld) > 1014 {
		return packet{}, fmt.Errorf("the payload is too big")
	}
	return packet{pType, seq, peer, port, payld}, nil
}

func (pkt packet) Bytes() []byte {
	buf := new(bytes.Buffer)
	buf.WriteByte(pkt.pType)
	buf.Write(pkt.seq)
	buf.Write(pkt.peer)
	buf.Write(pkt.port)
	buf.Write(pkt.payld)
	return buf.Bytes()
}

func toBig(x *byte) error {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, x)
	if err != nil {
		return err
	}
	return binary.Read(buf, binary.BigEndian, x)
}

func toBigEnd(x *[]byte) error {
	buf := new(bytes.Buffer)
	xb := make([]byte, len(*x))
	for i := range *x {
		xb[i] = (*x)[i]
	}
	err := binary.Write(buf, binary.BigEndian, &xb)
	if err != nil {
		return err
	}
	err = binary.Read(buf, binary.BigEndian, &xb)
	if err != nil {
		return err
	}
	x = &xb
	return nil
}
