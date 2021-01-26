package packets

import (
	"crypto/ed25519"
	"encoding/binary"
)

const (
	// ConnACKConnectionAccepted Connection accepted
	ConnACKConnectionAccepted = iota

	// ConnACKUnsupportedProtocol The Server does not support the level of the LightMQ protocol requested by the Client
	ConnACKUnsupportedProtocol

	// ConnACKServerUnavailable The Network Connection has been made but the LightMQ service is unavailable
	ConnACKServerUnavailable

	// ConnACKMalformedPayload Malformed payload
	ConnACKMalformedPayload

	// ConnACKUnauthorized The Client is not authorized to connect
	ConnACKUnauthorized
)

// ConnACK Packet is the packet sent by the Server in response to a CONNECT Packet received from a Client. The first packet sent from the Server to the Client MUST be a CONNACK Packet
type ConnACK struct {
	// e.g ConnACKConnectionAccepted
	ReturnCode byte
}

// Bytes converts ConnACK to bytes
func (c ConnACK) Bytes(skey ed25519.PrivateKey) []byte {
	payloadLength := make([]byte, 2)
	binary.BigEndian.PutUint16(payloadLength, 1)

	b := make([]byte, 68)
	b[0] = byte(TypeConnACK)
	b[1] = payloadLength[0]
	b[2] = payloadLength[1]
	b[3] = c.ReturnCode

	return b
}
