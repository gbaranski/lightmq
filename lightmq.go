package lightmq

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"sync"

	"github.com/gbaranski/lightmq/handlers"
	"github.com/gbaranski/lightmq/pkg/packets"
	"github.com/gbaranski/lightmq/pkg/types"
	"github.com/gbaranski/lightmq/pkg/utils"
	log "github.com/sirupsen/logrus"
)

// ClientList ...
type ClientList struct {
	s  []types.Client
	mu sync.RWMutex
}

// Add adds new client
func (l *ClientList) Add(c types.Client) {
	l.mu.Lock()
	l.s = append(l.s, c)
	l.mu.Unlock()
}

// Broker ...
type Broker struct {
	opts       Options
	clientList *ClientList
}

// New ...
func New(opts Options) (Broker, error) {
	broker := Broker{
		opts:       opts.Parse(),
		clientList: &ClientList{},
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

	fh, err := packets.ReadFixedHeader(conn)
	if err != nil {
		log.WithError(err).Error("fail read fixed header")
		return
	}

	if fh.ControlPacketType() != packets.TypeConnect {
		log.WithField("type", fh.ControlPacketType()).Info("Connection must start with CONNECT packet")
		return
	}

	len, err := utils.ReadLength(conn)
	if err != nil {
		log.WithError(err).Error("fail read leangth %s", err.Error())
		return
	}

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
	go b.clientList.Add(client)
	loge := log.WithFields(log.Fields{
		"clientID": client.ClientID,
		"username": client.Username,
		"ip":       client.IPAddress.String(),
	})

	loge.Info("Started connection")

	err = b.readLoop(conn, client)
	if err != nil {
		loge.WithError(err).Error("fail readLoop()")
	}
}

func (b Broker) readLoop(conn net.Conn, client types.Client) error {
	r := bufio.NewReader(conn)
	for {
		fh, err := packets.ReadFixedHeader(r)
		if err != nil {
			return fmt.Errorf("fail read fixed header %s", err.Error())
		}

		var handler handlers.Handler

		switch fh.ControlPacketType() {
		case packets.TypeSubscribe:
			fmt.Println("Subscribe not implemented yet")
			continue
		case packets.TypePingREQ:
			handler = handlers.OnPingReq
		case packets.TypePublish:
			handler = handlers.OnPublish
		default:
			return fmt.Errorf("unrecognized control packet type %x", fh.ControlPacketType())
		}

		len, err := utils.ReadLength(r)
		if err != nil {
			return fmt.Errorf("fail read len %s", err.Error())
		}

		err = handler(handlers.Packet{
			Client:          client,
			FixedHeader:     fh,
			RemainingLength: len,
			Reader:          r,
			Writer:          conn,
		})
		if err != nil {
			return fmt.Errorf("fail handle %x: %s", fh.ControlPacketType(), err.Error())
		}
	}

}
