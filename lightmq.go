package lightmq

import (
	"bytes"
	"fmt"
	"net"

	"github.com/gbaranski/lightmq/handlers"
	"github.com/gbaranski/lightmq/pkg/packets"
	log "github.com/sirupsen/logrus"
)

// Broker ...
type Broker struct {
	opts Options
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
	fixedHeader := make([]byte, 2)
	conn.Read(fixedHeader)

	var handler handlers.Handler

	controlPacketType := fixedHeader[0] >> 4
	switch controlPacketType {
	case packets.TypeConnect:
		handler = handlers.OnConnect
	default:
		log.WithField("type", fmt.Sprintf("%0b", controlPacketType)).Error("Unrecognized control packet type")
		return
	}

	// Optimize this size
	data := make([]byte, 65535)
	len, err := conn.Read(data)
	if err != nil {
		log.WithField("error", err.Error()).Error("Fail read data")
		return
	}

	c := handlers.Connection{
		Reader: bytes.NewReader(data),
		Len:    len,
	}
	res, err := handler(c)
	if err != nil {
		log.Error(err)
	}
	fmt.Println("res: ", res)
}
