package broker

import (
	"net"

	"github.com/gbaranski/lightmq/pkg/packets"
	log "github.com/sirupsen/logrus"
)

func (b *Broker) onConnect(conn net.Conn) (c Client, err error) {
	log.WithField("ip", conn.RemoteAddr().String()).Info("Starting new connection")

	cp, err := packets.ReadConnectPayload(conn)
	if err != nil {
		cack := packets.ConnACK{
			ReturnCode: packets.ConnACKMalformedPayload,
		}.Bytes(b.cfg.PrivateKey)
		conn.Write(cack[:])
		return c, err
	}

	cack := packets.ConnACK{
		ReturnCode: packets.ConnACKConnectionAccepted,
	}.Bytes(b.cfg.PrivateKey)

	_, err = conn.Write(cack[:])

	if err != nil {
		return c, err
	}

	return Client{
		ID:        cp.ClientID,
		IPAddress: conn.RemoteAddr(),
	}, nil
}
