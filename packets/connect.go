package packets

import (
	"bytes"
	"encoding/binary"
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
	Length      uint16
	ProtoName   [4]byte
	ProtoLevel  byte
	Flags       ConnectPacketFlags
	KeepAlive   uint16
	ClientID    []byte
	WillTopic   []byte
	WillMessage []byte
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
func ExtractConnectPacket(b []byte) (cp ConnectPacket, err error) {
	cp.FixedHeader = b[0]
	if cp.FixedHeader != ConnectFixedHeader {
		return cp, fmt.Errorf("invalid fixed header: %x", cp.FixedHeader)
	}
	cp.Length = binary.BigEndian.Uint16([]byte{b[1], b[2]})

	copy(cp.ProtoName[:], b[4:8])
	if !bytes.Equal(cp.ProtoName[:], []byte{0x4D, 0x51, 0x54, 0x54}) {
		return cp, fmt.Errorf("invalid proto name: %v", cp.ProtoName)
	}
	cp.ProtoLevel = b[8]
	if cp.ProtoLevel != SupportedProtoLevel {
		return cp, fmt.Errorf("unsupported proto level: %x", cp.ProtoLevel)
	}

	fb := b[9]
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
	cp.KeepAlive = binary.BigEndian.Uint16([]byte{b[11], b[12]})

	p := bytes.NewReader(b[13:])
	clientIDSize, err := p.ReadByte()
	if err != nil {
		return cp, fmt.Errorf("fail read clientID len %s", err.Error())
	}
	cp.ClientID = make([]byte, clientIDSize+1)
	_, err = p.Read(cp.ClientID)
	if err != nil {
		return cp, fmt.Errorf("fail read clientID %s", err.Error())
	}
	fmt.Println("clientID:", string(cp.ClientID))

	if cp.Flags.Will {
		willTopicSize, err := p.ReadByte()
		if err != nil {
			return cp, fmt.Errorf("fail read willTopicSize %s", err.Error())
		}
		cp.WillTopic = make([]byte, willTopicSize+1)
		_, err = p.Read(cp.WillTopic)
		if err != nil {
			return cp, fmt.Errorf("fail read willTopic %s", err.Error())
		}

		willMessageSize, err := p.ReadByte()
		if err != nil {
			return cp, fmt.Errorf("fail read willMessageSize %s", err.Error())
		}

		cp.WillMessage = make([]byte, willMessageSize+1)
		_, err = p.Read(cp.WillMessage)
		if err != nil {
			return cp, fmt.Errorf("fail read willMessage %s", err.Error())
		}
	}


	fmt.Printf("flags: %+v\n", cp.Flags)
	return
}
