package utils

import (
	"bytes"
	"encoding/binary"
	"io"
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

// ReadLength ...
func ReadLength(r io.Reader) (uint32, error) {
	len := uint32(0)

	buf := make([]byte, 1)
	buf[0] = 0b10000000

	for (buf[0] >> 7) == 1 {
		buf = make([]byte, 1)
		_, err := r.Read(buf)
		if err != nil {
			return 0, err
		}
		len += uint32(buf[0] & 0b01111111)
	}

	return len, nil
}
