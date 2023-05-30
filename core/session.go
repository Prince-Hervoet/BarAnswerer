package core

import (
	"ShareMemTCP/memory"
	"net"
	"time"
)

type Session struct {
	key            int64
	address        *Address
	connection     net.Conn
	mapping        *memory.ShareMemory
	data           map[string]any
	waitQueue      chan []byte
	startTimestamp int64
}

type InitSession struct {
	addr            *net.UDPAddr
	HandshakeStatus int8
}

func NewSession(key int64, address *Address, mapping *memory.ShareMemory, connection net.Conn) *Session {
	return &Session{
		key:            key,
		address:        address,
		connection:     connection,
		mapping:        mapping,
		data:           make(map[string]any),
		waitQueue:      make(chan []byte, 64),
		startTimestamp: time.Now().UnixMilli(),
	}
}
