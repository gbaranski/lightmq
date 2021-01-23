package handlers

import (
	"io"

	"github.com/gbaranski/lightmq/pkg/packets"
	"github.com/gbaranski/lightmq/pkg/types"
)

// Packet ...
type Packet struct {
	Client          types.Client
	FixedHeader     packets.FixedHeader
	RemainingLength uint32
	io.ByteReader
	io.Reader
	io.Writer
}

// Handler is type for packet handling function
type Handler = func(p Packet) error
