package handlers

import (
	"io"

	"github.com/gbaranski/lightmq/pkg/types"
)

// Handler is type for packet handling function
type Handler = func(reader io.ByteReader, w io.Writer, client types.Client) error
