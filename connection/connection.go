package connection

import (
	"net"
)

type Connection struct {
	OpenTimestamp int64
	ConnectionId  string
	Conn          net.Conn
}
