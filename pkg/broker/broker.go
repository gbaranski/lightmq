package broker

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"

	"github.com/gbaranski/lightmq/pkg/packets"
	"github.com/gbaranski/lightmq/pkg/utils"
	log "github.com/sirupsen/logrus"
)

// Client ...
type Client struct {
	ClientID  string
	IPAddress net.Addr
}

type packet struct {
	Client          Client
	RemainingLength uint32
	io.ByteReader
	io.Reader
	io.Writer
}

// Topic ...
type Topic struct {
	Name        string
	Subscribers map[string]struct{}
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

// Listen starts listening to incoming requests, this function is blocking
func (b *Broker) Listen() error {
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", b.cfg.Hostname, b.cfg.Port))
	if err != nil {
		return err
	}
	log.WithFields(log.Fields{
		"hostname": b.cfg.Hostname,
		"port":     b.cfg.Port,
	}).Info("Listening for incoming LightMQ connections")
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			return fmt.Errorf("fail accepting connection %s", err.Error())
		}
		go b.handleConnection(conn)
	}
}

func (b *Broker) handleConnection(conn net.Conn) {
	defer conn.Close()

	ptype, err := packets.ReadPacketType(conn)
	if err != nil {
		log.WithError(err).Error("fail read packet type")
		return
	}

	if ptype != packets.TypeConnect {
		log.WithField("type", ptype).Error("Connection must start with CONNECT packet")
		return
	}
	_, payload, err := packets.ReadSignedPayload(conn)
	if err != nil {
		log.WithError(err).Error("fail reading payload")
		return
	}

	client, err := b.onConnect(bytes.NewReader(payload), conn)
	if err != nil {
		log.WithError(err).Error("fail handle connection")
		return
	}
	go b.ClientStore.Add(client)
	loge := log.WithFields(log.Fields{
		"clientID": client.ClientID,
		"ip":       client.IPAddress.String(),
	})

	loge.Info("Started connection")

	err = b.readLoop(conn, client)
	if err != nil {
		loge.WithError(err).Error("fail readLoop()")
	}
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
			return fmt.Errorf("TypeSend not implemented")
		default:
			return fmt.Errorf("unrecognized control packet type %x", ptype)
		}

		len, err := utils.ReadLength(r)
		if err != nil {
			return fmt.Errorf("fail read len %s", err.Error())
		}

		err = handler(packet{
			Client:          client,
			RemainingLength: len,
			ByteReader:      r,
			Reader:          r,
			Writer:          conn,
		})
		if err != nil {
			return fmt.Errorf("fail handle %x: %s", ptype, err.Error())
		}
	}

}
