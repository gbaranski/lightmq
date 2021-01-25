package packets

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/gbaranski/lightmq/pkg/utils"
)

// SendPayload is payload of SEND packet
type SendPayload struct {
	ID    uint16
	Flags byte
	Data  []byte
}

// ReadSendPayload reads SendPayload from io.Reader
func ReadSendPayload(r io.Reader, plen uint16) (sp SendPayload, err error) {
	idBytes := make([]byte, 2)
	_, err = r.Read(idBytes)
	if err != nil {
		return sp, fmt.Errorf("fail read ID bytes %s", err.Error())
	}
	sp.ID = binary.BigEndian.Uint16(idBytes)

	sp.Flags, err = utils.ReadByte(r)
	if err != nil {
		return sp, fmt.Errorf("fail read flags byte %s", err.Error())
	}

	sp.Data = make([]byte, plen-3)
	_, err = r.Read(sp.Data)
	if err != nil {
		return sp, fmt.Errorf("fail read data %s", err.Error())
	}

	return sp, nil
}
