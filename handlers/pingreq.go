package handlers

import (
	"io"

	"github.com/gbaranski/lightmq/pkg/packets"
	"github.com/gbaranski/lightmq/pkg/types"
	log "github.com/sirupsen/logrus"
)

// OnPingReq ...
func OnPingReq(r io.ByteReader, w io.Writer, client types.Client) error {
	log.Info("Received PINGREQ")
	res := packets.PingRESP{}.Bytes()
	_, err := w.Write(res[:])
	if err != nil {
		return err
	}

	log.Info("Sent PINGRESP")
	return nil
}
