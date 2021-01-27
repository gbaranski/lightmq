package client

import (
	"fmt"
	"math"
	"math/rand"
	"net"

	"github.com/gbaranski/lightmq/pkg/packets"
	log "github.com/sirupsen/logrus"
)

// PacketChannels is struct which contains of channels for packets
type PacketChannels struct {
	ConnACK chan packets.ConnACKPayload
}

// Client ...
type Client struct {
	cfg     Config
	conn    net.Conn
	Packets PacketChannels
}

// New creates new client
func New(cfg Config) Client {
	return Client{
		cfg: cfg,
		Packets: PacketChannels{
			ConnACK: make(chan packets.ConnACKPayload),
		},
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
		OpCode:  packets.OpCodeConnect,
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

// Handler is type for packet handling function
type handler = func() error

// ReadLoop reads all data from connection in loop
func (c Client) ReadLoop() error {
	for {
		opcode, err := packets.ReadOpCode(c.conn)
		if err != nil {
			return fmt.Errorf("fail read packet type %s", err)
		}
		var handler handler

		switch opcode {
		case packets.OpCodeConnACK:
			handler = c.onConnACK
		case packets.OpCodePing:
			handler = c.onPing
		case packets.OpCodePong:
			handler = c.onPong
		default:
			return fmt.Errorf("no handler for opcode:%x", opcode)
		}

		log.WithField("opcode", fmt.Sprintf("0x%x", opcode)).Info("Handling packet")
		err = handler()
		if err != nil {
			return fmt.Errorf("fail handle %x: %s", opcode, err.Error())
		}
	}

}

// Send sends a data
func (c Client) Send(data []byte) error {
	payload := packets.SendPayload{
		ID:    uint16(rand.Intn(math.MaxInt16)),
		Flags: 0,
		Data:  data,
	}.Bytes()
	packet, err := packets.Packet{
		OpCode:  packets.OpCodeSend,
		Payload: payload,
	}.Bytes()
	if err != nil {
		return fmt.Errorf("fail encode payload %s", err.Error())
	}

	_, err = c.conn.Write(packet)

	return err
}

// Ping sends PING packet, returns ID of ping
func (c Client) Ping() (uint16, error) {
	payload := packets.PingPayload{
		ID: uint16(rand.Intn(math.MaxUint16)),
	}

	packet, err := packets.Packet{
		OpCode:  packets.OpCodePing,
		Payload: payload.Bytes(),
	}.Bytes()
	if err != nil {
		return 0, fmt.Errorf("fail encode packet")
	}
	_, err = c.conn.Write(packet)

	return payload.ID, err
}
