package packets

import (
	"bytes"
	"math"
	"math/rand"
	"testing"
)

func TestPongPayload(t *testing.T) {
	payload := PongPayload{
		ID: uint16(rand.Intn(math.MaxUint16)),
	}
	b := make([]byte, 0)
	b = append(b, 0) // MSB
	b = append(b, 2) // LSB
	b = append(b, payload.Bytes()...)

	readenPayload, err := ReadPongPayload(bytes.NewReader(b))
	if err != nil {
		t.Fatalf(err.Error())
	}
	if readenPayload.ID != payload.ID {
		t.Fatalf("Unexpected ID: %d, expected: %d", readenPayload.ID, payload.ID)
	}
}
