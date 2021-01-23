package packets

import (
	"fmt"
	"io"

	"github.com/gbaranski/lightmq/pkg/utils"
)

// PublishFixedHeader ...
type PublishFixedHeader byte

// ControlPacketType ...
func (h PublishFixedHeader) ControlPacketType() byte {
	return byte(h >> 4)
}

// Dup ...
func (h PublishFixedHeader) Dup() bool {
	return h&0b00001000 == 8
}

// QoS ...
func (h PublishFixedHeader) QoS() byte {
	return byte((h >> 1) & 0b00000011)
}

// Retain ...
func (h PublishFixedHeader) Retain() bool {
	return h&0b00000001 == 1
}

// PublishVariableHeader is variable header of Publish packet
type PublishVariableHeader struct {
	TopicName        []byte
	PacketIdentifier uint16
	Length           uint32
}

// ReadPublishVariableHeader reads PublishVariableHeader from io.Reader
func ReadPublishVariableHeader(b io.Reader) (h PublishVariableHeader, err error) {
	topicNameLength, err := utils.Read16BitInteger(b)
	if err != nil {
		return h, fmt.Errorf("fail read topic name length %s", err.Error())
	}
	h.TopicName = make([]byte, topicNameLength)
	_, err = b.Read(h.TopicName)
	if err != nil {
		return h, fmt.Errorf("fail read topicname %s", err.Error())
	}

	h.PacketIdentifier, err = utils.Read16BitInteger(b)
	if err != nil {
		return h, err
	}
	// 0,1 - uint16 topic name length
	// 1,2,3,4...n - topic name
	// n+1, n+2 - uint16 Packet Identifier
	h.Length = 2 + uint32(topicNameLength) + 2

	return h, nil
}
