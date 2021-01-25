package broker

import (
	"bufio"
	"fmt"
	"io"
	"net"

	"github.com/gbaranski/lightmq/pkg/packets"
)

// Client ...
type Client struct {
	ID        string
	IPAddress net.Addr
}

type packet struct {
	Client Client
	io.Writer
	io.Reader
}

// Handler is type for packet handling function
type handler = func(p packet) error

// Broker ...
type Broker struct {
	cfg         Config
	ClientStore *ClientStore
}

// New ...
func New(cfg Config) (Broker, error) {
	broker := Broker{
		cfg:         cfg.Parse(),
		ClientStore: NewClientStore(),
	}
	return broker, nil
}

func (b *Broker) readLoop(conn net.Conn, client Client) error {
	r := bufio.NewReader(conn)
	for {
		ptype, err := packets.ReadPacketType(r)
		if err != nil {
			return fmt.Errorf("fail read packet type %s", err.Error())
		}

		var handler handler

		switch ptype {
		case packets.TypeConnect:
			return fmt.Errorf("unexpected connect packet")
		case packets.TypeSend:
			handler = b.onSend
		default:
			return fmt.Errorf("unrecognized control packet type %x", ptype)
		}

		err = handler(packet{
			Client: client,
			Writer: conn,
			Reader: conn,
		})
		if err != nil {
			return fmt.Errorf("fail handle %x: %s", ptype, err.Error())
		}
	}

}