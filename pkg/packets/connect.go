package packets

import (
	"fmt"
	"io"
	"math"

	"github.com/gbaranski/lightmq/pkg/utils"
)

const (
	// ConnectPayloadChallengeSize size of the challage on connect packet
	ConnectPayloadChallengeSize int8 = 8
)

// ConnectPayload ...
type ConnectPayload struct {
	ClientID  string
	Challenge []byte
}

// ReadConnectPayload ...
func ReadConnectPayload(r io.Reader) (cp ConnectPayload, err error) {
	clientIDSize, err := utils.ReadByte(r)
	if err != nil {
		return cp, fmt.Errorf("fail read clientID len %s", err.Error())
	}
	clientID := make([]byte, clientIDSize)
	n, err := r.Read(clientID)
	if err != nil {
		return cp, fmt.Errorf("fail read clientID %s", err.Error())
	}
	if n != int(clientIDSize) {
		return cp, fmt.Errorf("invalid clientID size")
	}
	cp.ClientID = string(clientID)

	cp.Challenge = make([]byte, ConnectPayloadChallengeSize)
	n, err = r.Read(cp.Challenge)
	if err != nil {
		return cp, fmt.Errorf("fail read challenge %s", err.Error())
	}
	if n != int(ConnectPayloadChallengeSize) {
		return cp, fmt.Errorf("invalid challange size")
	}

	return cp, nil
}

// Bytes converts ConnectPayload to bytes
func (cp *ConnectPayload) Bytes() ([]byte, error) {
	if len(cp.ClientID) > math.MaxUint8 {
		return nil, fmt.Errorf("ClientID is too big")
	}
	// 1 			  <- ClientID Size byte
	// len(clientID)  <- ClientID length
	// len(challenge) <- Challenge length

	// Think about optimizing this one
	p := make([]byte, 0)
	p = append(p, uint8(len(cp.ClientID)))
	p = append(p, []byte(cp.ClientID)...)
	p = append(p, cp.Challenge...)

	return p, nil

}
