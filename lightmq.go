package lightmq

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"

	"github.com/gbaranski/lightmq/handlers"
	"github.com/gbaranski/lightmq/pkg/packets"
	"github.com/gbaranski/lightmq/pkg/types"
	"github.com/gbaranski/lightmq/pkg/utils"
	log "github.com/sirupsen/logrus"
)

// Broker ...
type Broker struct {
	opts    Options
	clients []types.Client
}

// New ...
func New(opts Options) (Broker, error) {
	broker := Broker{
		opts: opts.Parse(),
	}
	return broker, nil
}

// Listen starts listening to incoming requests, this function is blocking
func (b Broker) Listen() error {
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", b.opts.Hostname, b.opts.Port))
	if err != nil {
		return err
	}
	log.WithFields(log.Fields{
		"hostname": b.opts.Hostname,
		"port":     b.opts.Port,
	}).Info("Listening for incoming MQTT connections")
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			return fmt.Errorf("fail accepting connection %s", err.Error())
		}
		go b.handleConnection(conn)

	}
}

func (b Broker) handleConnection(conn net.Conn) {
	defer conn.Close()

	fixedHeader := make([]byte, 1)
	_, err := conn.Read(fixedHeader)
	if err != nil {
		log.WithError(err).Error("fail read fixed header")
		return
	}

	controlPacketType := fixedHeader[0] >> 4
	if controlPacketType != packets.TypeConnect {
		log.WithField("type", controlPacketType).Info("Connection must start with CONNECT packet")
		return
	}

	len, err := utils.ReadLength(conn)
	if err != nil {
		log.WithError(err).Error("fail read leangth %s", err.Error())
		return
	}

	// Optimize this size
	data := make([]byte, len)
	_, err = conn.Read(data)
	if err != nil {
		log.WithField("error", err.Error()).Error("Fail read data")
		return
	}
	client, err := handlers.OnConnect(bytes.NewReader(data), conn)
	if err != nil {
		log.Error("fail connection %s", err.Error())
		return
	}
	b.clients = append(b.clients, client)
	loge := log.WithFields(log.Fields{
		"clientID": client.ClientID,
		"username": client.Username,
		"ip":       client.IPAddress.String(),
	})

	loge.Info("Started connection")

	err = b.readLoop(conn, client)
	if err != nil {
		loge.WithError(err).Error("fail loop conn data")
	}
}

func (b Broker) readLoop(conn net.Conn, client types.Client) error {
	r := bufio.NewReader(conn)
	for {
		fh, err := r.ReadByte()
		if err != nil {
			if err == io.EOF {
				continue
			}
			return fmt.Errorf("fail read fixed header %s", err.Error())
		}
		if fh == 0 {
			fmt.Println("Skipping 0")
			continue
		}

		controlPacketType := fh >> 4

		var handler handlers.Handler

		switch controlPacketType {
		case packets.TypeSubscribe:
			fmt.Println("Subscribe not implemented yet")
			continue
		case packets.TypePingREQ:
			handler = handlers.OnPingReq
		default:
			fmt.Printf("fh: %0b\n", fh)
			return fmt.Errorf("unrecognized control packet type %x", controlPacketType)
		}

		len, err := utils.ReadLength(r)
		if err != nil {
			return fmt.Errorf("fail read len %s", err.Error())
		}
		// Optimize this size
		data := make([]byte, len)
		_, err = r.Read(data)
		if err != nil {
			return err
		}

		err = handler(r, conn, client)
		if err != nil {
			return fmt.Errorf("fail handle %x: %s", controlPacketType, err.Error())
		}
	}

}
