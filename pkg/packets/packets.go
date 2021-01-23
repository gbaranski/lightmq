package packets

import (
	"bytes"
	"errors"
	"fmt"
	"io"
)

const (
	// TypeConnect - Client request to connect to Server
	//
	// Direction: Client to Server
	TypeConnect byte = iota + 1

	// TypeConnACK - Connect acknowledgment
	//
	// Direction: Server to Client
	TypeConnACK

	// TypePublish - Publish message
	//
	// Direction: Server to Client or Client to Server
	TypePublish

	// TypePubACK - Publish acknowledgment
	//
	// Direction: Server to Client or Client to Server
	TypePubACK

	// TypePubREC - Publish received (assured delivery part 1)
	//
	// Direction: Server to Client or Client to Server
	TypePubREC

	// TypePubREL - Publish release (assured delivery part 2)
	//
	// Direction: Server to Client or Client to Server
	TypePubREL

	// TypePubCOMP - Publish complete (assured delivery part 3)
	//
	// Direction: Server to Client or Client to Server
	TypePubCOMP

	// TypeSubscribe - Client subscribe request
	//
	// Direction: Client to Server
	TypeSubscribe

	// TypeSubACK - Subscribe acknowledgment
	//
	// Direction: Server to Client
	TypeSubACK

	// TypeUnsubscribe - Unsubscribe request
	//
	// Direction: Client to Server
	TypeUnsubscribe

	// TypeUnsubACK - Unsubscribe acknowledgment
	//
	// Direction: Server to Client
	TypeUnsubACK

	// TypePingREQ - Ping request
	//
	// Direction: Client to Server
	TypePingREQ

	// TypePingRESP - Ping response
	//
	// Direction: Server to Client
	TypePingRESP

	// TypeDisconnect - Client is disconnecting
	//
	// Direction: Client to Server
	TypeDisconnect
)

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

var (
	// ErrUnacceptableProtocol ...
	ErrUnacceptableProtocol = errors.New("unnaceptable protocol")
)

// VerifyProtoName verifies if proto name is equal to ['M','Q','T','T']
func VerifyProtoName(r *bytes.Reader) error {
	var expb = [6]byte{0x0, 0x4, 0x4D, 0x51, 0x54, 0x54}
	for i, exp := range expb {
		b, err := r.ReadByte()
		if err != nil {
			return fmt.Errorf("can't read byte at %d: %s", i, err.Error())
		}
		if b != byte(exp) {
			return fmt.Errorf("invalid byte at %d: exp: %x, rec: %x", i, exp, b)
		}
	}
	return nil
}

// FixedHeader ...
type FixedHeader byte

// ControlPacketType ...
func (h FixedHeader) ControlPacketType() byte {
	return byte(h >> 4)
}

// ReadFixedHeader reads fixed header from io.Reader
func ReadFixedHeader(r io.Reader) (FixedHeader, error) {
	b := make([]byte, 1)
	n, err := r.Read(b)
	if err != nil {
		return 0, err
	}
	if n != 1 {
		return 0, fmt.Errorf("read invalid amount, exp: 1, n: %d", n)
	}
	return FixedHeader(b[0]), nil
}
