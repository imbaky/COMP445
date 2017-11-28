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

type Packet struct {
	PType byte
	Seq   []byte // 4bytes
	Peer  []byte // 4bytes
	Port  []byte // 2bytes
	Payld []byte // max 1014 bytes
}

func (p *Packet) GetSequence() uint32 {
	return binary.BigEndian.Uint32(p.Seq)
}
func MakePacket(pType byte, seq, peer, port, payld []byte) (Packet, error) {
	if pType != ACK && pType != NACK && pType != SYN && pType != SYNACK {
		return Packet{}, fmt.Errorf("packet type must be one of the following %v %v %v %v", ACK, NACK, SYN, SYNACK)
	}
	if len(seq) != 4 {
		return Packet{}, fmt.Errorf("seq is not 4 bytes")
	}

	fmt.Printf("seqstart %x\n", seq)
	seq, _ = toBigEnd4(seq)
	fmt.Printf("seqend %x\n", seq)
	if len(peer) != 4 {
		return Packet{}, fmt.Errorf("Peer is not 4 bytes")
	}
	peer, _ = toBigEnd4(peer)
	if len(port) != 2 {
		return Packet{}, fmt.Errorf("Port is not 2 bytes")
	}
	port, _ = toBigEnd2(port)
	if len(payld) > 1014 {
		return Packet{}, fmt.Errorf("the payload is too big")
	}
	return Packet{pType, seq, peer, port, payld}, nil
}

func (pkt Packet) Bytes() []byte {
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

func FromBytes(raw []byte) (Packet, error) {
	if len(raw) < MinLen || len(raw) > MaxLen {
		return Packet{}, fmt.Errorf("packet is too big or too small")
	}
	var pkt Packet
	pkt.PType = raw[0]
	pkt.Seq = raw[1:5]
	pkt.Peer = raw[5:9]
	pkt.Port = raw[9:11]
	pkt.Payld = raw[13:]
	return pkt, nil
}
