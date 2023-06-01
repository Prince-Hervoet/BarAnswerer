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
	SidMap     map[int]string
	sessions   map[string]*Session
	callBacks  map[string]func(buf []byte)
	EpollIndex *EpollInfo
	selection  int8
}

// 创建一个服务端分享者
func newServerSharer() *ServerSharer {
	return &ServerSharer{
		SidMap:    make(map[int]string),
		sessions:  make(map[string]*Session),
		callBacks: make(map[string]func(buf []byte)),
		selection: util.SERVER,
	}
}

// 创建一个seesion并设置对应的map
func (here *ServerSharer) PushSessionMap(sessionId string, port int, mapping *memory.ShareMemory, connection net.Conn, fd int) {
	t := NewSession(sessionId, port, mapping, connection)
	if connection == nil {
		t.fd = fd
	}
	here.SidMap[t.fd] = sessionId
	here.sessions[sessionId] = t
}

// 删除结构体中的映射
func (here *ServerSharer) DeleteSession(fd int) {
	delete(here.sessions, here.SidMap[fd])
	delete(here.SidMap, fd)
}

// 关闭对话后的释放资源操作
func (here *ServerSharer) RecoverResource(sessionId string) {
	session := here.sessions[sessionId]
	if session.mapping != nil {
		here.sessions[sessionId].mapping.Close() //移除共享内存映射
	}
	here.DeleteSession(here.sessions[sessionId].fd) //移除session映射
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

		if here.EpollIndex != nil {
			here.EpollIndex.AddEvent(int32(fd.Fd()), here.MemShareTcpDeal)
		}
	}
}

// 关闭会话
func (here *ServerSharer) Close(sessionId string) error {
	if _, has := here.sessions[sessionId]; !has {
		return errors.New("sessionId error")
	}
	if here.EpollIndex != nil {
		here.EpollIndex.DeleteEvent(int32(here.sessions[sessionId].fd))
	}
	here.RecoverResource(sessionId)
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

func (here *ServerSharer) Read(buffer []byte, sessionId string) (int32, error) {
	if _, has := here.sessions[sessionId]; !has {
		return 0, errors.New("sessionId error")
	}
	session := here.sessions[sessionId]
	count, err := session.mapping.Read(buffer)
	if err != nil {
		return 0, err
	}
	return count, nil
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
		here.waveHandDeal(buf, fd)
	case util.NOTICE_MESSAGE:
		here.noticeDeal(buf, fd)
	default:
		fmt.Println("server memory share tcp error")
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

	here.PushSessionMap(id, 0, sm, nil, fd)

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

// epoll中处理写通知的服务端函数
func (here *ServerSharer) noticeDeal(buf []byte, fd int) {

	sessionId := here.SidMap[fd]
	if _, has := here.callBacks[sessionId]; has {
		here.callBacks[sessionId](buf) //调用用户设置的回调函数
	}

	here.sessions[sessionId].mapping.ChangeStatus(0)
	message := CreateMessage(util.CLIENT, util.NOTICE_MESSAGE, nil)
	unix.Write(fd, message)
}

// epoll中处理写通知的服务端函数
func (here *ServerSharer) waveHandDeal(buf []byte, fd int) {
	var sessionId string
	data, PayLoadLen, _ := ReadMessege(buf, util.SHACK_HAND_MESSAGE, util.SERVER)

	for i := 0; i < int(PayLoadLen); i++ {
		if data[i] == '\n' {
			sessionId = string(data[0:i])
		}
	}

	if here.EpollIndex != nil {
		fmt.Print("fd:")
		fmt.Println(here.sessions[sessionId].fd)
		here.EpollIndex.DeleteEvent(int32(here.sessions[sessionId].fd))
	}

	here.RecoverResource(sessionId)

}
