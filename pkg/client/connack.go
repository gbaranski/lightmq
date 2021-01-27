package client

import (
	"fmt"

	"github.com/gbaranski/lightmq/pkg/packets"
)

func (c Client) onConnACK() error {
	payload, err := packets.ReadConnACKPayload(c.conn)
	if err != nil {
		return fmt.Errorf("fail parse payload %s", err.Error())
	}
	c.Packets.ConnACK <- payload

	return nil
}
