package client

import (
	"fmt"
	"math"
	"math/rand"
	"net"

	"github.com/gbaranski/lightmq/pkg/packets"
)

// Client ...
type Client struct {
	cfg  Config
	conn net.Conn
}

// New creates new client
func New(cfg Config) Client {
	return Client{
		cfg: cfg,
	}
}

// Connect connects to the specified host with specified port
func (c *Client) Connect() error {
	var err error
	c.conn, err = net.Dial("tcp", fmt.Sprintf("%s:%d", c.cfg.Hostname, c.cfg.Port))
	if err != nil {
		return err
	}
	payload, err := packets.ConnectPayload{
		ClientID: c.cfg.ClientID,
		// Fill it up later
		Challenge: []byte{0, 0, 0, 0, 0, 0, 0, 0},
	}.Bytes()
	if err != nil {
		return fmt.Errorf("fail convert connect payload to bytes %s", err.Error())
	}

	p, err := packets.Packet{
		Type:    packets.TypeConnect,
		Payload: payload,
	}.Bytes()
	if err != nil {
		return fmt.Errorf("fail convert connect packet to bytes %s", err.Error())
	}
	_, err = c.conn.Write(p)
	if err != nil {
		return fmt.Errorf("fail write CONNECT packet %s", err.Error())
	}

	return nil
}

// Send sends a data
func (c Client) Send(data []byte) error {
	payload := packets.SendPayload{
		ID:    uint16(rand.Intn(math.MaxInt16)),
		Flags: 0,
		Data:  data,
	}.Bytes()
	packet, err := packets.Packet{
		Type:    packets.TypeSend,
		Payload: payload,
	}.Bytes()
	if err != nil {
		return fmt.Errorf("fail encode payload %s", err.Error())
	}

	_, err = c.conn.Write(packet)

	return err
}
