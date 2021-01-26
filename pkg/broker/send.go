package broker

import (
	"fmt"

	"github.com/gbaranski/lightmq/pkg/packets"
	log "github.com/sirupsen/logrus"
)

// OnSend should be executed on SEND packet
func (b *Broker) onSend(p packet) error {
	_, err := packets.ReadSignature(p)
	if err != nil {
		return fmt.Errorf("fail read sig %s", err.Error())
	}

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
