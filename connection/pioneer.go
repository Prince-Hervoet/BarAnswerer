package connection

import (
	"ShareMemTCP/util"
	"fmt"
	"net"
	"time"
)

type Pioneer struct {
	connections []*Connection
}

// 打开连接
func (here *Pioneer) OpenConnection(address string) (int64, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println("error connecting")
		return -1, err
	}
	id := util.NextId()
	nc := &Connection{
		OpenTimestamp: time.Now().UnixMilli(),
		ConnectionId:  id,
		Conn:          conn,
	}
	here.connections = append(here.connections, nc)
	return id, nil
}

// 检查连接
func (here *Pioneer) CheckConnection() {

}

// 关闭连接
func (here *Pioneer) CloseConnection() {

}
