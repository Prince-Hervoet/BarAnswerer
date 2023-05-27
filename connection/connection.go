package connection

import "net"

type Connection struct {
	OpenTimestamp int64
	ConnectionId  int64
	Conn          net.Conn
}
