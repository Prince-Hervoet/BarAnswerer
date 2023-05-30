package core

import (
	"ShareMemTCP/memory"
	"ShareMemTCP/util"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"golang.org/x/sys/unix"
)

type ServerSharer struct {
	sessions  map[string]*Session
	Epoll     *EpollInfo
	selection int8
}

// 创建一个服务端分享者
func NewServerSharer() *ServerSharer {
	return &ServerSharer{
		sessions:  make(map[string]*Session),
		Epoll:     &EpollInfo{},
		selection: util.SERVER,
	}
}

func (here *ServerSharer) SetCallback() {

}

// 监听线程tcp连接线程
func (here *ServerSharer) Listen(port int) {
	// 监听端口
	portStr := strconv.FormatInt(int64(port), 10)
	ln, err := net.Listen("tcp", ":"+portStr)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ln.Close()

	file, err := ln.(*net.TCPListener).File()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("TCP server is listening on port %d, listen fd:%d\n", port, file.Fd())

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		fd, err := conn.(*net.TCPConn).File()
		fmt.Printf("listen fd:%d\n", fd.Fd())
		if err != nil {
			continue
		}

		here.Epoll.AddEvent(int32(fd.Fd()), here.MemShareTcpDeal)
	}
}

// 关闭会话
func (here *ServerSharer) Close(sessionId string) error {
	if _, has := here.sessions[sessionId]; !has {
		return errors.New("sessionId error")
	}
	here.Epoll.DeleteEvent(int32(here.sessions[sessionId].fd))
	return nil
}

// 当客户端发来请求时的回复函数
func (here *ServerSharer) Reply(content []byte, sessionId string) error {
	if _, has := here.sessions[sessionId]; !has {
		return errors.New("sessionId error")
	}
	session := here.sessions[sessionId]
	session.mapping.Write(content)
	return nil
}

func (here *ServerSharer) Read() {

}

// 开启监听端口，等待握手事件到来然后开启共享内存
func (here *ServerSharer) Open(port int) {

}

// 客户端发送报文处理回调函数
func (here *ServerSharer) MemShareTcpDeal(buf []byte, fd int) {

	side, MessageType, err := ReadMessegeHeader(buf)
	if err != nil {
		fmt.Println("read header error")
		return
	}

	if side != util.CLIENT {
		return
	}

	switch MessageType {
	case util.SHACK_HAND_MESSAGE:
		here.shakeDeal(buf, fd)
	case util.WAVE_HAND_MESSAGE:

	case util.NOTICE_MESSAGE:

	default:

	}

}

// epoll中处理握手的服务器函数
func (here *ServerSharer) shakeDeal(buf []byte, fd int) {
	data, dataLen, _ := ReadMessege(buf, util.SHACK_HAND_MESSAGE, util.CLIENT)
	var cap int
	if dataLen > 2 {
		cap = int(util.BytesToInt32(data))
	} else {
		cap = int(util.BytesToInt16(data))
	}

	sm := memory.OpenShareMemory()
	fileName, _ := gonanoid.New()

	id, _ := gonanoid.New()
	t := NewSession(id, 0, sm, nil)
	t.fd = fd
	here.sessions[id] = t

	filePath, err := sm.OpenFile(fileName, int32(cap))
	if err != nil {
		fmt.Println(err)
		here.Close(id)
		return
	}

	str := &strings.Builder{}
	str.WriteString(id)
	str.WriteByte(byte('\n'))
	str.WriteString(filePath)
	str.WriteByte(byte('\n'))
	payLoad := []byte(str.String())

	sendBuf := CreateMessage(util.SERVER, util.SHACK_HAND_MESSAGE, payLoad)

	unix.Write(fd, sendBuf)
}
