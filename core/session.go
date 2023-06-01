package core

import (
	"ShareMemTCP/memory"
	"net"
	"time"
)

type Session struct {
	sessionId      string
	port           int
	connection     net.Conn
	fd             int
	mapping        *memory.ShareMemory
	data           map[string]any
	startTimestamp int64
}

func NewSession(sessionId string, port int, mapping *memory.ShareMemory, connection net.Conn) *Session {
	ans := &Session{
		sessionId:      sessionId,
		port:           port,
		connection:     connection,
		mapping:        mapping,
		data:           make(map[string]any),
		startTimestamp: time.Now().UnixMilli(),
	}
	if connection != nil {
		file, _ := connection.(*net.TCPConn).File()
		ans.fd = int(file.Fd())
	}
	return ans
}
