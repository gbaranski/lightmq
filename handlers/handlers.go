package handlers

import (
	"bytes"
)

// Handler is type for packet handling function
type Handler = func(Connection) ([]byte, error)

// Connection ...
type Connection struct {
	Reader *bytes.Reader
	Len    int
}
