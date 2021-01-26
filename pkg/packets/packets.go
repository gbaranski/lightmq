package packets

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"

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
	// SignatureSize is size of signature in bytes
	SignatureSize uint8 = 64
)

// Payload is type for LightMQ Payload
type Payload []byte

// PacketType defines type of packet, can be one of Type...
type PacketType byte

// ReadPacketType reads packet type and returns it
func ReadPacketType(r io.Reader) (PacketType, error) {
	b, err := utils.ReadByte(r)
	return PacketType(b), err
}

// Signature is ed25519 signature of the payload
type Signature []byte

// ReadSignature reads signature, that must be called after reading Packet Type, but before reading Payload length
func ReadSignature(r io.Reader) (sig Signature, err error) {
	sig = make(Signature, SignatureSize)
	n, err := r.Read(sig)
	if err != nil {
		return sig, fmt.Errorf("fail read signature %s", err.Error())
	}
	if n != int(SignatureSize) {
		return sig, fmt.Errorf("invalid signature length: %d", n)
	}

	return sig, nil
}

// ReadPayloadSize reads size length, that must be called after reading signature but before reading payload
func ReadPayloadSize(r io.Reader) (uint16, error) {
	return utils.Read16BitInteger(r)
}

// Packet is type for LightMQ Packet
type Packet struct {
	Type      PacketType
	Signature Signature
	Payload   Payload
}

// Bytes converts packet to bytes which can be directly sent
func (p Packet) Bytes() ([]byte, error) {
	if len(p.Payload) > math.MaxUint16 {
		return nil, fmt.Errorf("payload is too big: %d", len(p.Payload))
	}
	plen := make([]byte, 2)
	binary.BigEndian.PutUint16(plen, uint16(len(p.Payload)))
	// Optimize it later
	b := make([]byte, 0)
	b = append(b, byte(p.Type))
	b = append(b, p.Signature...)
	b = append(b, plen...)
	b = append(b, p.Payload...)
	return b, nil
}
