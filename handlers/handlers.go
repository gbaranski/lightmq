package handlers

import (
	"bytes"
	"io"
)

// Handler is type for packet handling function
type Handler = func(Connection) error

// Connection ...
type Connection struct {
	Reader *bytes.Reader
	Len    int
	io.Writer
}
