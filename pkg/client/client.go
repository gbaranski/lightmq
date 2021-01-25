package client

import (
	"crypto/ed25519"
	"fmt"
	"net"

	"github.com/gbaranski/lightmq/pkg/packets"
)

// Client ...
type Client struct {
	cfg Config
}

// New creates new client
func New(cfg Config) Client {
	return Client{
		cfg: cfg,
	}
}

// Connect connects to the specified host with specified port
func (c *Client) Connect() error {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", c.cfg.Hostname, c.cfg.Port))
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

	sig := ed25519.Sign(c.cfg.PrivateKey, payload)
	p, err := packets.Packet{
		Type:      packets.TypeConnect,
		Signature: sig,
		Payload:   payload,
	}.Bytes()
	if err != nil {
		return fmt.Errorf("fail convert connect packet to bytes %s", err.Error())
	}
	_, err = conn.Write(p)
	if err != nil {
		return fmt.Errorf("fail write CONNECT packet %s", err.Error())
	}

	return nil
}
