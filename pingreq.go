package lightmq

import (
	"github.com/gbaranski/lightmq/pkg/packets"
	log "github.com/sirupsen/logrus"
)

func (b *Broker) onPingReq(p packet) error {
	log.Info("Received PINGREQ")
	res := packets.PingRESP{}.Bytes()
	_, err := p.Write(res[:])
	if err != nil {
		return err
	}

	log.Info("Sent PINGRESP")
	return nil
}
