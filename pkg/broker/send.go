package broker

import (
	"github.com/gbaranski/lightmq/pkg/packets"
	log "github.com/sirupsen/logrus"
)

// OnSend should be executed on SEND packet
func (b *Broker) onSend(p packet) error {
	psize, err := packets.ReadPayloadSize(p)
	if err != nil {
		return err
	}
	sp, err := packets.ReadSendPayload(p, psize)
	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"clientID": p.Client.ID,
		"msgID":    sp.ID,
		"data":     string(sp.Data),
	}).Info("Received SEND packet")

	return nil
}
