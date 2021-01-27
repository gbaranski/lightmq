package packets

import (
	"io"

	"github.com/gbaranski/lightmq/pkg/utils"
)

// PongPayload is payload for PONG Control packet
type PongPayload struct {
	ID uint16
}

// ReadPongPayload reads pong payload from io.Reader
func ReadPongPayload(r io.Reader) (PingPayload, error) {
	id, err := utils.Read16BitInteger(r)
	if err != nil {
		return PingPayload{}, err
	}
	return PingPayload{
		ID: id,
	}, nil
}

// Bytes convert PongPayload to bytes
func (p PongPayload) Bytes() (b []byte) {
	b = make([]byte, 2)
	b[0] = byte(p.ID >> 8)
	b[1] = byte(p.ID)

	return b
}
