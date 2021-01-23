package types

import "net"

// Client ...
type Client struct {
	ClientID  string
	Username  string
	IPAddress net.Addr
	KeepAlive uint16
}
