package handlers

import (
	"fmt"

	"github.com/gbaranski/lightmq/packets"
	log "github.com/sirupsen/logrus"
)

// OnConnect should be executed when CONNECT packet is received
func OnConnect(c Connection) ([]byte, error) {
	cp, err := packets.ReadConnectPacket(c.Reader)
	if err != nil {
		return nil, err
	}

	log.WithFields(log.Fields{
		"clientID":    string(cp.Payload.ClientID),
		"username":    string(cp.Payload.Username),
		"password":    string(cp.Payload.Password),
		"willMessage": string(cp.Payload.WillMessage),
		"willTopic":   string(cp.Payload.WillTopic),
		"flags":       fmt.Sprintf("%+v", cp.Flags),
	}).Info("Received Connect packet")

	// cack := packets.ConnACK{
	// 	Flags: packets.ConnACKFlags{
	// 		SessionPresent: false,
	// 	},
	// 	ReturnCode: packets.ConnACKConnectionAccepted,
	// }.Bytes()
	// conn.Write(cack[:])
	// log.Info("Sent ConnACK packet")
	return nil, nil

}
