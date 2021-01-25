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
	b := make([]byte, 68)

	sig := ed25519.Sign(skey, []byte{c.ReturnCode})
	payloadLength := make([]byte, 2)
	binary.BigEndian.PutUint16(payloadLength, 1)

	b = append(b, byte(TypeConnACK)) // 1 byte
	b = append(b, sig...)            // 64 bytes
	b = append(b, payloadLength...)  // 2 bytes
	b = append(b, c.ReturnCode)      // 1 byte

	return b
}
