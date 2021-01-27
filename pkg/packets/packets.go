package packets

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"

	"github.com/gbaranski/lightmq/pkg/utils"
)

const (
	// OpCodeConnect - Client request to connect to Server operation code
	//
	// Direction: Client to Server
	OpCodeConnect OpCode = iota + 1

	// OpCodeConnACK - Connect acknowledgment operation code
	//
	// Direction: Server to Client
	OpCodeConnACK

	// OpCodeSend - Send message operation code
	//
	// Direction: Server to Client or Client to Server
	OpCodeSend

	// OpCodeSendRESP - Send Response operation code
	//
	// Direction: Server to Client or Client to Server
	OpCodeSendRESP

	// OpCodePing - Ping request operation code
	//
	// Direction: Server to Client or Client to Server
	OpCodePing

	// OpCodePong - Ping acknowledgmenet operation code
	//
	// Direction: Server to Client or Client to Server
	OpCodePong
)

// Payload is type for LightMQ Payload
type Payload []byte

// OpCode defnes opcode of packet
type OpCode byte

// ReadOpCode reads packet type and returns it
func ReadOpCode(r io.Reader) (OpCode, error) {
	b, err := utils.ReadByte(r)
	return OpCode(b), err
}

// ReadPayloadSize reads size length, that must be called before reading payload
func ReadPayloadSize(r io.Reader) (uint16, error) {
	return utils.Read16BitInteger(r)
}

// Packet is type for LightMQ Packet
type Packet struct {
	OpCode  OpCode
	Payload Payload
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
	b = append(b, byte(p.OpCode))
	b = append(b, plen...)
	b = append(b, p.Payload...)
	return b, nil
}
