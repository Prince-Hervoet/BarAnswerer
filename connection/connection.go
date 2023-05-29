package connection

import (
	"net"

	"golang.org/x/sys/unix"
)

type Connection struct {
	OpenTimestamp int64
	ConnectionId  int64
	Conn          net.Conn
}

type EpollInfo struct {
	EpollFd  int
	EventNum int
	events   []unix.EpollEvent
	mp       map[int32]func([]byte ,int)
}
