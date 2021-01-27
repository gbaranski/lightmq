package client

import (
	"fmt"

	"github.com/gbaranski/lightmq/pkg/packets"
	log "github.com/sirupsen/logrus"
)

func (c Client) onPing() error {
	pingp, err := packets.ReadPingPayload(c.conn)
	if err != nil {
		return fmt.Errorf("fail parse payload %s", err.Error())
	}
	log.WithFields(log.Fields{
		"pingID": pingp.ID,
	}).Info("Received PING packet")
	pongp := packets.PongPayload{

		ID: pingp.ID,
	}
	_, err = c.conn.Write(pongp.Bytes())
	if err != nil {
		return fmt.Errorf("fail snd pong %s", err.Error())
	}

	log.WithFields(log.Fields{
		"pongID": pongp.ID,
	}).Info("Sent PONG packet")

	return nil
}
