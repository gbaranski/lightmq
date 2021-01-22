package packets

import (
	"bytes"
	"fmt"
)

// ConnectPacketFlags ...
type ConnectPacketFlags struct {
	CleanSession bool
	Will         bool
	WillQoS      uint8
	WillRetain   bool
	Username     bool
	Password     bool
}

// ConnectPacket ...
type ConnectPacket struct {
	FixedHeader byte
	ProtoName   [6]byte
	ProtoLevel  byte
	Flags       ConnectPacketFlags
}

const (
	// AtMostOnceQoS says that message will be delivered max 1 time
	AtMostOnceQoS uint8 = 0x00

	// AtLeastOnceQoS says that message will be delivered at least once
	AtLeastOnceQoS uint8 = 0x01

	// ExactlyOnceQoS says that message will be delivered exactly once
	ExactlyOnceQoS uint8 = 0x02

	// ConnectFixedHeader ...
	ConnectFixedHeader uint8 = 0x10

	// SupportedProtoLevel defines the supported level by the broker
	//
	// It might be changed later to array of uint8
	SupportedProtoLevel uint8 = 0x4
)

// ExtractConnectPacket ...
func ExtractConnectPacket(p []byte) (cp ConnectPacket, err error) {
	cp.FixedHeader = p[0]
	if cp.FixedHeader != ConnectFixedHeader {
		return cp, fmt.Errorf("invalid fixed header: %x", cp.FixedHeader)
	}
	copy(cp.ProtoName[:], p[3:9])
	if !bytes.Equal(cp.ProtoName[:], []byte{0x0, 0x4, 0x4D, 0x51, 0x54, 0x54}) {
		return cp, fmt.Errorf("invalid proto name: %v", cp.ProtoName)
	}
	cp.ProtoLevel = p[9]
	if cp.ProtoLevel != SupportedProtoLevel {
		return cp, fmt.Errorf("unsupported proto level: %x", cp.ProtoLevel)
	}

	fb := p[10]
	fmt.Printf("10:%0b\n", p[10])
	if fb&0b1 != 0x0 {
		return cp, fmt.Errorf("invalid reserved bit")
	}
	cp.Flags = ConnectPacketFlags{
		CleanSession: (fb>>1)&1 == 1,
		Will:         (fb>>2)&1 == 1,
		WillQoS:      (fb >> 3) & 0b11,
		WillRetain:   (fb>>5)&1 == 1,
		Username:     (fb>>7)&1 == 1,
		Password:     (fb>>6)&1 == 1,
	}

	fmt.Printf("flags: %+v\n", cp.Flags)
	return
}
