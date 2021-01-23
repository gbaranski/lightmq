package packets

import (
	"bytes"
	"fmt"
)

const (
	// TypeConnect - Client request to connect to Server
	//
	// Direction: Client to Server
	TypeConnect byte = iota + 1

	// TypeConnACK - Connect acknowledgment
	//
	// Direction: Server to Client
	TypeConnACK

	// TypePublish - Publish message
	//
	// Direction: Server to Client or Client to Server
	TypePublish

	// TypePubACK - Publish acknowledgment
	//
	// Direction: Server to Client or Client to Server
	TypePubACK

	// TypePubREC - Publish received (assured delivery part 1)
	//
	// Direction: Server to Client or Client to Server
	TypePubREC

	// TypePubREL - Publish release (assured delivery part 2)
	//
	// Direction: Server to Client or Client to Server
	TypePubREL

	// TypePubCOMP - Publish complete (assured delivery part 3)
	//
	// Direction: Server to Client or Client to Server
	TypePubCOMP

	// TypeSubscribe - Client subscribe request
	//
	// Direction: Client to Server
	TypeSubscribe

	// TypeSubACK - Subscribe acknowledgment
	//
	// Direction: Server to Client
	TypeSubACK

	// TypeUnsubscribe - Unsubscribe request
	//
	// Direction: Client to Server
	TypeUnsubscribe

	// TypeUnsubACK - Unsubscribe acknowledgment
	//
	// Direction: Server to Client
	TypeUnsubACK

	// TypePingREQ - Ping request
	//
	// Direction: Client to Server
	TypePingREQ

	// TypePingRESP - Ping response
	//
	// Direction: Server to Client
	TypePingRESP

	// TypeDisconnect - Client is disconnecting
	//
	// Direction: Client to Server
	TypeDisconnect
)

// VerifyProtoName verifies if proto name is equal to ['M','Q','T','T']
func VerifyProtoName(r *bytes.Reader) error {
	var expb = [6]byte{0x0, 0x4, 0x4D, 0x51, 0x54, 0x54}
	for i, exp := range expb {
		b, err := r.ReadByte()
		if err != nil {
			return fmt.Errorf("fail read at %d: %s", i, err.Error())
		}
		if b != byte(exp) {
			return fmt.Errorf("invalid byte at %d: exp: %x, rec: %x", i, exp, b)
		}
	}
	return nil
}
