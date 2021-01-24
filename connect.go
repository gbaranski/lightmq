package lightmq

import (
	"bytes"
	"net"

	"github.com/gbaranski/lightmq/pkg/packets"
	log "github.com/sirupsen/logrus"
)

func (b *Broker) onConnect(r *bytes.Reader, conn net.Conn) (client Client, err error) {
	client.IPAddress = conn.RemoteAddr()
	log.WithField("ip", client.IPAddress.String()).Info("Starting new connection")

	h, err := packets.ReadConnectVariableHeader(r)
	if err != nil {
		if err == packets.ErrUnacceptableProtocol {
			cack := packets.ConnACK{
				Flags: packets.ConnACKFlags{
					SessionPresent: false,
				},
				ReturnCode: packets.ConnACKConnectionAccepted,
			}.Bytes()
			conn.Write(cack[:])
			return client, err
		}
		return client, err
	}
	client.KeepAlive = h.KeepAlive

	payload, err := packets.ReadConnectPayload(r, h.Flags)
	if err != nil {
		return Client{}, err
	}
	client.ClientID = string(payload.ClientID)
	client.Username = string(payload.Username)

	cack := packets.ConnACK{
		Flags: packets.ConnACKFlags{
			SessionPresent: false,
		},
		ReturnCode: packets.ConnACKConnectionAccepted,
	}.Bytes()
	_, err = conn.Write(cack[:])
	if err != nil {
		return Client{}, err
	}

	return client, nil
}
