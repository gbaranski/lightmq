package broker

import (
	"fmt"

	"github.com/gbaranski/lightmq/pkg/packets"
	log "github.com/sirupsen/logrus"
)

// onPing should be executed on PING packet
func (b *Broker) onPing(p packet) error {
	pingp, err := packets.ReadPingPayload(p)
	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"clientID": p.Client.ID,
		"pingID":   pingp.ID,
	}).Info("Received PING packet")
	pongp := packets.PongPayload{
		ID: pingp.ID,
	}
	_, err = p.Write(pongp.Bytes())
	if err != nil {
		return fmt.Errorf("fail snd pong %s", err.Error())
	}

	log.WithFields(log.Fields{
		"clientID": p.Client.ID,
		"pongID":   pongp.ID,
	}).Info("Sent PONG packet")

	return err
}
