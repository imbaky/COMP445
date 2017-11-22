package packet

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const (
	MinLen int  = 11
	MaxLen int  = 1024
	ACK    byte = 0xFF
	NACK   byte = 0x00
	SYN    byte = 0xF0
	SYNACK byte = 0x0F
)

type packet struct {
	PType byte
	Seq   []byte // 4bytes
	Peer  []byte // 4bytes
	Port  []byte // 2bytes
	Payld []byte // max 1014 bytes
}

func Packet(pType byte, seq, peer, port, payld []byte) (packet, error) {
	if pType != ACK && pType != NACK && pType != SYN && pType != SYNACK {
		return packet{}, fmt.Errorf("packet type must be one of the following %v %v %v %v", ACK, NACK, SYN, SYNACK)
	}
	if len(seq) != 4 {
		return packet{}, fmt.Errorf("seq is not 4 bytes")
	}

	fmt.Printf("seqstart %x\n", seq)
	seq, _ = toBigEnd4(seq)
	fmt.Printf("seqend %x\n", seq)
	if len(peer) != 4 {
		return packet{}, fmt.Errorf("Peer is not 4 bytes")
	}
	peer, _ = toBigEnd4(peer)
	if len(port) != 2 {
		return packet{}, fmt.Errorf("Port is not 2 bytes")
	}
	port, _ = toBigEnd2(port)
	if len(payld) > 1014 {
		return packet{}, fmt.Errorf("the payload is too big")
	}
	return packet{pType, seq, peer, port, payld}, nil
}

func (pkt packet) Bytes() []byte {
	buf := new(bytes.Buffer)
	buf.WriteByte(pkt.PType)
	buf.Write(pkt.Seq)
	buf.Write(pkt.Peer)
	buf.Write(pkt.Port)
	buf.Write(pkt.Payld)
	return buf.Bytes()
}

func toBigEnd2(x []byte) ([]byte, error) {
	fmt.Printf("start %x\n", x)
	b := make([]byte, 2)
	littleEnd := binary.LittleEndian.Uint16(x)
	binary.BigEndian.PutUint16(b, littleEnd)
	fmt.Printf("end %x\n", b)
	return b, nil
}

func toBigEnd4(x []byte) ([]byte, error) {
	fmt.Printf("start %x\n", x)
	b := make([]byte, 4)
	littleEnd := binary.LittleEndian.Uint32(x)
	binary.BigEndian.PutUint32(b, littleEnd)
	fmt.Printf("end %x\n", b)
	return b, nil
}

func fromBytes(raw []byte) (packet, error) {
	if len(raw) < MinLen || len(raw) > MaxLen {
		return packet{}, fmt.Errorf("packet is too big or too small")
	}
	var pkt packet
	pkt.PType = raw[0]
	pkt.Seq = raw[1:5]
	pkt.Peer = raw[5:9]
	pkt.Port = raw[9:12]
	pkt.Payld = raw[12:]
	return pkt, nil
}
