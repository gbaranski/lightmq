package packets

const (
	// ConnACKConnectionAccepted Connection accepted
	ConnACKConnectionAccepted = iota
	// ConnACKUnacceptableProtocol The Server does not support the level of the MQTT protocol requested by the Client
	ConnACKUnacceptableProtocol

	// ConnACKIdentifierRejected The Client identifier is correct UTF-8 but not allowed by the Server
	ConnACKIdentifierRejected

	// ConnACKServerUnavailable The Network Connection has been made but the MQTT service is unavailable
	ConnACKServerUnavailable

	// ConnACKBadUsernameOrPassword The data in the user name or password is malformed
	ConnACKBadUsernameOrPassword

	// ConnACKNotAuthorized The Client is not authorized to connect
	ConnACKNotAuthorized
)

// ConnACKFlags ...
type ConnACKFlags struct {
	SessionPresent bool
}

// ConnACK Packet is the packet sent by the Server in response to a CONNECT Packet received from a Client. The first packet sent from the Server to the Client MUST be a CONNACK Packet [MQTT-3.2.0-1].
type ConnACK struct {
	Flags ConnACKFlags

	// e.g ConnACKConnectionAccepted
	ReturnCode byte
}

// Bytes converts ConnACK to bytes
func (c ConnACK) Bytes() (b [3]byte) {
	b[0] = TypeConnACK << 4 // <- 0b00100000 Fixed Header

	if c.Flags.SessionPresent {
		b[1] = 0b00000001
	} else {
		b[1] = 0b00000000
	}

	b[2] = c.ReturnCode

	return
}
