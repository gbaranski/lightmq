package broker

import (
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
	for {
		ptype, err := packets.ReadOpCode(conn)
		if err != nil {
			return fmt.Errorf("fail read packet type %s", err.Error())
		}

		var handler handler

		switch ptype {
		case packets.OpCodeConnect:
			return fmt.Errorf("unexpected connect packet")
		case packets.OpCodeSend:
			handler = b.onSend
		case packets.OpCodePing:
			handler = b.onPing
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
