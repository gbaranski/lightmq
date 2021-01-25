package broker

import (
	"bytes"
	"net"

	"github.com/gbaranski/lightmq/pkg/packets"
	log "github.com/sirupsen/logrus"
)

func (b *Broker) onConnect(p *bytes.Reader, conn net.Conn) (Client, error) {
	log.WithField("ip", conn.RemoteAddr().String()).Info("Starting new connection")

	cp, err := packets.ReadConnectPayload(p)

	cack := packets.ConnACK{
		ReturnCode: packets.ConnACKConnectionAccepted,
	}.Bytes(b.cfg.PrivateKey)
	_, err = conn.Write(cack[:])
	if err != nil {
		return Client{}, err
	}

	return Client{
		ClientID:  cp.ClientID,
		IPAddress: conn.RemoteAddr(),
	}, nil
}
