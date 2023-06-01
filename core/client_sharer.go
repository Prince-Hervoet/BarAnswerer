package core

import (
	"ShareMemTCP/memory"
	"ShareMemTCP/util"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"golang.org/x/sys/unix"
)

type ClientSharer struct {
	SidMap     map[int]string
	sessions   map[string]*Session
	EpollIndex *EpollInfo
	selection  int8
}

// 创建客户端分享者
func newClientSharer() *ClientSharer {
	return &ClientSharer{
		SidMap:    make(map[int]string),
		sessions:  make(map[string]*Session),
		selection: util.CLIENT,
	}
}

// 挂载在epoll上的客户端读取事件,在通知中会取消在epoll中的挂载
func (here *ClientSharer) ClientMemTcpCallBack(buf []byte, Fd int) {

}

// 对话关闭后回收资源的函数
func (here *ClientSharer) RecoverResource(sessionId string) {
	here.sessions[sessionId].mapping.Close()        //移除共享内存映射
	here.DeleteSession(here.sessions[sessionId].fd) //移除session映射
}

// 创建一个seesion并设置对应的map
func (here *ClientSharer) PushSessionMap(sessionId string, port int, mapping *memory.ShareMemory, connection net.Conn, fd int) {
	t := NewSession(sessionId, port, mapping, connection)
	if connection == nil {
		t.fd = fd
	}
	here.SidMap[t.fd] = sessionId
	here.sessions[sessionId] = t
}

func (here *ClientSharer) DeleteSession(fd int) {
	delete(here.sessions, here.SidMap[fd])
	delete(here.SidMap, fd)
}

// 链接对端进程，请求开启共享内存（传入的对端进程的地址、需要的内存大小、连接的标识）
func (here *ClientSharer) Link(port int, cap int32) (string, error) {
	//tcp连接建立
	portStr := strconv.FormatInt(int64(port), 10)
	conn, err1 := net.Dial("tcp", ":"+portStr)
	if err1 != nil {
		fmt.Println("error client connecting")
		conn.Close()
		return "", err1
	}

	fmt.Println("TCP connect")
	//共享内存协议握手
	sessionId, filePath, err2 := here.sheckHand(cap, conn)
	if err2 != nil {
		fmt.Println("error client shakeHand")
		conn.Close()
		return "", err2
	}
	fmt.Println("client hand shake over")

	shareMemory := memory.OpenShareMemory()
	err := shareMemory.LinkFile(filePath, cap)
	if err != nil {
		fmt.Println("shakeHand return fileName error")
		conn.Close()
		return "", err
	}

	//保存session并添加到epoll中
	here.PushSessionMap(sessionId, port, shareMemory, conn, 0)
	here.EpollIndex.AddEvent(int32(here.sessions[sessionId].fd), here.ClientMemTcpCallBack)

	return sessionId, nil
}

// 共享内存协议栈
func (here *ClientSharer) sheckHand(cap int32, conn net.Conn) (string, string, error) {
	// 初始化传输报文
	data := util.Int32ToBytes(cap)
	message := CreateMessage(util.CLIENT, util.SHACK_HAND_MESSAGE, data)

	//发送报文
	_, err := conn.Write(message)
	if err != nil {
		fmt.Println("error menTCP client write connecting")
		return "", "", err
	}

	//等待服务端初始化完成
	readMessege := make([]byte, 1024)
	_, err = conn.Read(readMessege)
	if err != nil {
		fmt.Println("error menTCP client read connecting")
		return "", "", err
	}

	var sessionId string
	var filePath string
	lastIndex := 0
	data, PayLoadLen, _ := ReadMessege(readMessege, util.SHACK_HAND_MESSAGE, util.SERVER)

	for i := 0; i < int(PayLoadLen); i++ {
		if data[i] == '\n' {
			if len(sessionId) > 0 {
				filePath = string(data[lastIndex:i])
			} else {
				sessionId = string(data[lastIndex:i])
				lastIndex = i + 1
			}
		}
	}
	fmt.Printf("clint get seesionID:%s filePath:%s\n", sessionId, filePath)
	return sessionId, filePath, nil
}

// 发送数据到共享内存中（需要通知对端进程进行读取）
func (here *ClientSharer) Send(data []byte, sessionId string) error {
	if _, has := here.sessions[sessionId]; !has {
		return errors.New("invalid sessionId")
	}
	session := here.sessions[sessionId]
	err := session.mapping.Write(data)
	if err != nil {
		return err
	}

	//向共享内存写数据
	session.mapping.Write(data)
	session.mapping.ChangeStatus(1)

	// 发送通知
	message := CreateMessage(util.CLIENT, util.NOTICE_MESSAGE, nil)
	session.connection.Write(message)

	// 等待服务端处理完毕,要先取下后再阻塞读取
	here.EpollIndex.DeleteEvent(int32(here.sessions[sessionId].fd))

	buf := make([]byte, 128)
	//设置读取的截止时间
	session.connection.SetReadDeadline(time.Now().Add(1 * time.Second))
	session.connection.Read(buf)
	_, _, err = ReadMessege(buf, util.NOTICE_MESSAGE, util.SERVER)
	if err != nil {
		here.DeleteSession(session.fd)
		return err
	}

	here.EpollIndex.AddEvent(int32(here.sessions[sessionId].fd), here.ClientMemTcpCallBack)

	return nil
}

// 关闭会话
func (here *ClientSharer) Close(sessionId string) error {
	if _, has := here.sessions[sessionId]; !has {
		return errors.New("sessionId error")
	}
	if here.EpollIndex != nil {
		here.EpollIndex.DeleteEvent(int32(here.sessions[sessionId].fd))
	}

	session := here.sessions[sessionId]

	//传输结束报文
	str := &strings.Builder{}
	str.WriteString(sessionId)
	str.WriteByte(byte('\n'))
	payLoad := []byte(str.String())

	sendBuf := CreateMessage(util.SERVER, util.SHACK_HAND_MESSAGE, payLoad)

	unix.Write(session.fd, sendBuf)

	here.RecoverResource(sessionId)
	return nil
}

// 读取共享内存中数据
// func (here *ClientSharer) Read(bs []byte, sessionId string) (int32, error) {
// 	if _, has := here.sessions[sessionId]; !has {
// 		return 0, errors.New("invalid sessionId")
// 	}
// 	session := here.sessions[sessionId]
// 	count, err := session.mapping.Read(bs)
// 	if err != nil {
// 		return 0, err
// 	}
// 	return count, nil
// }
