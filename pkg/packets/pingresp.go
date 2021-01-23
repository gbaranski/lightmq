package packets

// PingRESP ...
type PingRESP struct {
}

// Bytes converts PingRESP to byte array
func (p PingRESP) Bytes() (b [2]byte) {
	b[0] = TypePingRESP << 4 // <- 0b11010000 Fixed header
	b[1] = 0x0               // <- Remaining length(0)

	return b
}
