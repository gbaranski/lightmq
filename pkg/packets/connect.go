package packets

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// ConnectFlags ...
type ConnectFlags struct {
	CleanSession bool
	Will         bool
	WillQoS      uint8
	WillRetain   bool
	Username     bool
	Password     bool
}

// ConnectPayload ...
type ConnectPayload struct {
	ClientID    []byte
	WillTopic   []byte
	WillMessage []byte
	Username    []byte
	Password    []byte
}

// Connect ...
type Connect struct {
	Length     uint16
	ProtoName  [4]byte
	ProtoLevel byte
	Flags      ConnectFlags
	KeepAlive  uint16
	Payload    ConnectPayload
}

const (
	// AtMostOnceQoS says that message will be delivered max 1 time
	AtMostOnceQoS uint8 = 0x00

	// AtLeastOnceQoS says that message will be delivered at least once
	AtLeastOnceQoS uint8 = 0x01

	// ExactlyOnceQoS says that message will be delivered exactly once
	ExactlyOnceQoS uint8 = 0x02

	// SupportedProtoLevel defines the supported level by the broker
	//
	// It might be changed later to array of uint8
	SupportedProtoLevel uint8 = 0x4
)

// ReadConnectFlags ...
func ReadConnectFlags(p *bytes.Reader) (ConnectFlags, error) {
	b, err := p.ReadByte()
	if err != nil {
		return ConnectFlags{}, err
	}
	if b&0b1 != 0x0 {
		return ConnectFlags{}, fmt.Errorf("invalid reserved bit: %x", b&0b1)
	}

	cf := ConnectFlags{
		CleanSession: (b>>1)&1 == 1,
		Will:         (b>>2)&1 == 1,
		WillQoS:      (b >> 3) & 0b11,
		WillRetain:   (b>>5)&1 == 1,
		Username:     (b>>7)&1 == 1,
		Password:     (b>>6)&1 == 1,
	}
	if (!cf.Will && cf.WillQoS > 0) || (!cf.Will && cf.WillRetain) {
		return cf, fmt.Errorf("invalid flags, WillQoS/WillRetain can be true only if Will is true")
	}
	if cf.WillQoS != AtMostOnceQoS && cf.WillQoS != AtLeastOnceQoS && cf.WillQoS != ExactlyOnceQoS {
		return cf, fmt.Errorf("invalid QoS level %x", cf.WillQoS)
	}

	return cf, nil
}

// ReadConnectPayload ...
func ReadConnectPayload(p *bytes.Reader, f ConnectFlags) (cpp ConnectPayload, err error) {
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
