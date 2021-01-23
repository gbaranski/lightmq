package handlers

import (
	"github.com/gbaranski/lightmq/pkg/packets"
	log "github.com/sirupsen/logrus"
)

// OnPingReq ...
func OnPingReq(p Packet) error {
	log.Info("Received PINGREQ")
	res := packets.PingRESP{}.Bytes()
	_, err := p.Write(res[:])
	if err != nil {
		return err
	}

	log.Info("Sent PINGRESP")
	return nil
}
