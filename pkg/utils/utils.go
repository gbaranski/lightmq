package utils

import (
	"bytes"
	"encoding/binary"
)

// Read16BitInteger reads 16 bit integer from bytes reader
func Read16BitInteger(r *bytes.Reader) (uint16, error) {
	bytes := make([]byte, 2)
	_, err := r.Read(bytes)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint16(bytes), nil
}
