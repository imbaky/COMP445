package packet

import (
	"reflect"
	"testing"
)

func TestPacket(t *testing.T) {
	type args struct {
		pType PacketType
		seq   SequenceNumber
		peer  PeerAddr
		port  PortNumber
		payld Payload
	}
	tests := []struct {
		name string
		args args
		want packet
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Packet(tt.args.pType, tt.args.seq, tt.args.peer, tt.args.port, tt.args.payld); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Packet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_packet_Bytes(t *testing.T) {
	type fields struct {
		pType PacketType
		seq   SequenceNumber
		peer  PeerAddr
		port  PortNumber
		payld Payload
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pkt := packet{
				pType: tt.fields.pType,
				seq:   tt.fields.seq,
				peer:  tt.fields.peer,
				port:  tt.fields.port,
				payld: tt.fields.payld,
			}
			if got := pkt.Bytes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("packet.Bytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPortNumber_toBigEnd(t *testing.T) {
	tests := []struct {
		name    string
		x       *PortNumber
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.x.toBigEnd(); (err != nil) != tt.wantErr {
				t.Errorf("PortNumber.toBigEnd() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPeerAddr_toBigEnd(t *testing.T) {
	tests := []struct {
		name    string
		x       *PeerAddr
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.x.toBigEnd(); (err != nil) != tt.wantErr {
				t.Errorf("PeerAddr.toBigEnd() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSequenceNumber_toBigEnd(t *testing.T) {
	tests := []struct {
		name    string
		x       *SequenceNumber
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.x.toBigEnd(); (err != nil) != tt.wantErr {
				t.Errorf("SequenceNumber.toBigEnd() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPacketType_toBigEnd(t *testing.T) {
	tests := []struct {
		name    string
		p       *PacketType
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.p.toBigEnd(); (err != nil) != tt.wantErr {
				t.Errorf("PacketType.toBigEnd() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPacketType_bytes(t *testing.T) {
	tests := []struct {
		name string
		p    PacketType
		want []byte
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.bytes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PacketType.bytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSequenceNumber_bytes(t *testing.T) {
	tests := []struct {
		name string
		s    SequenceNumber
		want []byte
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.bytes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SequenceNumber.bytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPeerAddr_bytes(t *testing.T) {
	tests := []struct {
		name string
		p    PeerAddr
		want []byte
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.bytes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PeerAddr.bytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPortNumber_bytes(t *testing.T) {
	tests := []struct {
		name string
		p    PortNumber
		want []byte
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.bytes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PortNumber.bytes() = %v, want %v", got, tt.want)
			}
		})
	}
}
