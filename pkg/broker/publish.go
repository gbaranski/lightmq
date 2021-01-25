package broker

import (
	"fmt"

	"github.com/gbaranski/lightmq/pkg/packets"
	log "github.com/sirupsen/logrus"
)

// OnPublish should be executed on PUBLISH packet
func (b *Broker) onPublish(p packet) error {
	vh, err := packets.ReadPublishVariableHeader(p.Reader)
	if err != nil {
		return fmt.Errorf("fail read variable header %s", err.Error())
	}
	data := make([]byte, int(p.RemainingLength-vh.Length))
	_, err = p.Read(data)
	if err != nil {
		return fmt.Errorf("fail read data %s", err.Error())
	}
	log.WithFields(log.Fields{
		"topic":   string(vh.TopicName),
		"payload": string(data),
	}).Info("Received PUBLISH")

	return nil
}
