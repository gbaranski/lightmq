package handlers

import (
	"io"

	"github.com/gbaranski/lightmq/pkg/types"
)

// Packet ...
type Packet struct {
	Client types.Client
	io.Reader
	io.Writer
}

// Handler is type for packet handling function
type Handler = func(p Packet) error
