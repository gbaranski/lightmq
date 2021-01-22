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

// ConnectPacketPayload ...
type ConnectPacketPayload struct {
	ClientID    []byte
	WillTopic   []byte
	WillMessage []byte
	Username    []byte
	Password    []byte
}

// ConnectPacket ...
type ConnectPacket struct {
	FixedHeader byte
	Length      uint16
	ProtoName   [4]byte
	ProtoLevel  byte
	Flags       ConnectPacketFlags
	KeepAlive   uint16
	Payload     ConnectPacketPayload
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

func extractConnectFlags(b byte) ConnectPacketFlags {
	return ConnectPacketFlags{
		CleanSession: (b>>1)&1 == 1,
		Will:         (b>>2)&1 == 1,
		WillQoS:      (b >> 3) & 0b11,
		WillRetain:   (b>>5)&1 == 1,
		Username:     (b>>7)&1 == 1,
		Password:     (b>>6)&1 == 1,
	}
}

func extractConnectPayload(p *bytes.Reader, f ConnectPacketFlags) (cpp ConnectPacketPayload, err error) {
	clientIDSize, err := p.ReadByte()
	if err != nil {
		return cpp, fmt.Errorf("fail read clientID len %s", err.Error())
	}
	cpp.ClientID = make([]byte, clientIDSize)
	_, err = p.Read(cpp.ClientID)
	if err != nil {
		return cpp, fmt.Errorf("fail read clientID %s", err.Error())
	}

	p.ReadByte() // Consume null terminator

	if f.Will {
		willTopicSize, err := p.ReadByte()
		if err != nil {
			return cpp, fmt.Errorf("fail read willTopicSize %s", err.Error())
		}
		cpp.WillTopic = make([]byte, willTopicSize)
		_, err = p.Read(cpp.WillTopic)
		if err != nil {
			return cpp, fmt.Errorf("fail read willTopic %s", err.Error())
		}
		p.ReadByte() // Consume null terminator

		willMessageSize, err := p.ReadByte()
		if err != nil {
			return cpp, fmt.Errorf("fail read willMessageSize %s", err.Error())
		}

		cpp.WillMessage = make([]byte, willMessageSize)
		_, err = p.Read(cpp.WillMessage)
		if err != nil {
			return cpp, fmt.Errorf("fail read willMessage %s", err.Error())
		}

		p.ReadByte() // Consume null terminator
	}

	if f.Username {
		usernameSize, err := p.ReadByte()
		if err != nil {
			return cpp, fmt.Errorf("fail read usernameSize %s", err.Error())
		}
		cpp.Username = make([]byte, usernameSize)
		_, err = p.Read(cpp.Username)
		if err != nil {
			return cpp, fmt.Errorf("fail read username %s", err.Error())
		}
	}

	if f.Password {
		passwordSizeBytes := make([]byte, 2)
		if _, err = p.Read(passwordSizeBytes); err != nil {
			return cpp, fmt.Errorf("fail read passwordSizeBytes %s", err.Error())
		}
		passwordSize := binary.BigEndian.Uint16(passwordSizeBytes)
		cpp.Password = make([]byte, passwordSize)
		_, err = p.Read(cpp.Password)
		if err != nil {
			return cpp, fmt.Errorf("fail read password %s", err.Error())
		}
	}
	return cpp, nil
}

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
	cp.Flags = extractConnectFlags(fb)

	cp.KeepAlive = binary.BigEndian.Uint16([]byte{b[11], b[12]})
	p := bytes.NewReader(b[13:])
	cp.Payload, err = extractConnectPayload(p, cp.Flags)

	return cp, err
}
