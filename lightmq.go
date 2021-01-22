package lightmq

import (
	"fmt"
	"net"

	"github.com/gbaranski/lightmq/packets"
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
	buf := make([]byte, 1000)
	n, err := conn.Read(buf)
	if err != nil {
		log.Info("fail read data %s", err.Error())
		return
	}
	log.Info("read %d bytes", n)
	log.Info("bytes: %v", buf)
	_, err = packets.ExtractConnectPacket(buf)
	if err != nil {
		log.Error("Fail extract connect packet", err.Error())
	}
}
