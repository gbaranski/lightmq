package packets

import (
	"bytes"
	"fmt"
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
func ReadConnectPayload(p *bytes.Reader) (cp ConnectPayload, err error) {
	clientIDSize, err := p.ReadByte()
	if err != nil {
		return cp, fmt.Errorf("fail read clientID len %s", err.Error())
	}
	clientID := make([]byte, clientIDSize)
	n, err := p.Read(clientID)
	if err != nil {
		return cp, fmt.Errorf("fail read clientID %s", err.Error())
	}
	if n != int(clientIDSize) {
		return cp, fmt.Errorf("invalid clientID size")
	}
	cp.ClientID = string(clientID)

	cp.Challenge = make([]byte, ConnectPayloadChallengeSize)
	n, err = p.Read(cp.Challenge)
	if err != nil {
		return cp, fmt.Errorf("fail read challenge %s", err.Error())
	}
	if n != int(ConnectPayloadChallengeSize) {
		return cp, fmt.Errorf("invalid challange size")
	}

	return cp, nil
}
