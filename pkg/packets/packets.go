package packets

import (
	"fmt"
	"io"

	"github.com/gbaranski/lightmq/pkg/utils"
)

const (
	// TypeConnect - Client request to connect to Server
	//
	// Direction: Client to Server
	TypeConnect PacketType = iota + 1

	// TypeConnACK - Connect acknowledgment
	//
	// Direction: Server to Client
	TypeConnACK

	// TypeSend - Send message
	//
	// Direction: Server to Client or Client to Server
	TypeSend

	// TypeSendRESP - Send Response
	//
	// Direction: Server to Client or Client to Server
	TypeSendRESP
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

	// SignatureSize is size of signature in bytes
	SignatureSize uint8 = 64
)

// PacketType defines type of packet, can be one of Type...
type PacketType byte

// ReadPacketType reads packet type and returns it
func ReadPacketType(r io.Reader) (PacketType, error) {
	b, err := utils.ReadByte(r)
	return PacketType(b), err
}

// Payload is payload for the LightMQ Packet
type Payload []byte

// Signature is ed25519 signature of the payload
type Signature []byte

// ReadSignedPayload start with reading signature, then reads length of the data and the data	.
func ReadSignedPayload(r io.Reader) (sig Signature, p Payload, err error) {
	sig = make(Signature, SignatureSize)
	n, err := r.Read(sig)
	if err != nil {
		return sig, p, fmt.Errorf("fail read signature %s", err.Error())
	}
	if n != int(SignatureSize) {
		return sig, p, fmt.Errorf("invalid signature length: %d", n)
	}

	plen, err := utils.Read16BitInteger(r)
	if err != nil {
		return sig, p, fmt.Errorf("fail read payload len %s", err.Error())
	}
	p = make(Payload, plen)
	n, err = r.Read(p)
	if err != nil {
		return sig, p, fmt.Errorf("fail read payload %s", err.Error())
	}
	if uint16(n) != plen {
		return sig, p, fmt.Errorf("invalid payload length %s", err.Error())
	}
	return sig, p, nil
}
