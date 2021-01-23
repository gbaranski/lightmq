package handlers

import (
	"fmt"

	"github.com/gbaranski/lightmq/pkg/packets"
	"github.com/gbaranski/lightmq/pkg/utils"
	log "github.com/sirupsen/logrus"
)

// OnConnect should be executed when CONNECT packet is received
func OnConnect(c Connection) error {
	log.Info("New connection")
	if err := packets.VerifyProtoName(c.Reader); err != nil {
		return fmt.Errorf("fail verify proto name %s", err.Error())
	}
	protolevel, err := c.Reader.ReadByte()
	if err != nil {
		return fmt.Errorf("fail read proto level %s", err.Error())
	}
	if protolevel != packets.SupportedProtoLevel {
		cack := packets.ConnACK{
			ReturnCode: packets.ConnACKUnacceptableProtocol,
		}.Bytes()
		c.Write(cack[:])
		return fmt.Errorf("fail verify proto name %s", err.Error())
	}
	cf, err := packets.ReadConnectFlags(c.Reader)
	if err != nil {
		return err
	}

	keepAlive, err := utils.Read16BitInteger(c.Reader)
	if err != nil {
		return fmt.Errorf("fail read keepalive %s", err.Error())
	}
	c.Reader.ReadByte() // <- for some reason its required

	payload, err := packets.ReadConnectPayload(c.Reader, cf)
	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"clientID":    string(payload.ClientID),
		"username":    string(payload.Username),
		"password":    string(payload.Password),
		"willMessage": string(payload.WillMessage),
		"willTopic":   string(payload.WillTopic),
		"keepAlive":   keepAlive,
		"flags":       fmt.Sprintf("%+v", cf),
	}).Info("Received Connect packet")

	cack := packets.ConnACK{
		Flags: packets.ConnACKFlags{
			SessionPresent: false,
		},
		ReturnCode: packets.ConnACKConnectionAccepted,
	}.Bytes()
	_, err = c.Write(cack[:])
	if err != nil {
		return err
	}
	log.Info("Sent CONNACK")

	return err
}
