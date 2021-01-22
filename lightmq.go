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
	buf := make([]byte, 100)
	_, err := conn.Read(buf)
	if err != nil {
		log.Info("fail read data %s", err.Error())
		return
	}
	cp, err := packets.ExtractConnectPacket(buf)
	log.WithFields(log.Fields{
		"clientID":    string(cp.Payload.ClientID),
		"username":    string(cp.Payload.Username),
		"password":    string(cp.Payload.Password),
		"willMessage": string(cp.Payload.WillMessage),
		"willTopic":   string(cp.Payload.WillTopic),
		"flags":       fmt.Sprintf("%+v", cp.Flags),
	}).Info("Received connect packet")

	if err != nil {
		log.Error("Fail extract connect packet", err.Error())
	}
}
