package core

import (
	"ShareMemTCP/memory"
	"ShareMemTCP/util"
	"errors"
	"fmt"
	"net"
	"strconv"
)

type ClientSharer struct {
	sessions  map[string]*Session
	selection int8
}

// 创建客户端分享者
func NewClientSharer() *ClientSharer {
	return &ClientSharer{
		sessions:  make(map[string]*Session),
		selection: util.CLIENT,
	}
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

	//共享内存协议握手
	sessionId, filePath, err2 := here.sheckHand(cap, conn)
	if err2 != nil {
		fmt.Println("error client shakeHand")
		conn.Close()
		return "", err2
	}

	shareMemory := memory.OpenShareMemory()
	err := shareMemory.LinkFile(filePath, cap)
	if err != nil {
		fmt.Println("shakeHand return fileName error")
		conn.Close()
		return "", err
	}

	here.sessions[sessionId] = NewSession(sessionId, port, shareMemory, conn)
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
				lastIndex = i
				sessionId = string(data[lastIndex:i])
			}
		}
	}
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
	return nil
}

// 读取共享内存中数据
func (here *ClientSharer) Read(bs []byte, sessionId string) (int32, error) {
	if _, has := here.sessions[sessionId]; !has {
		return 0, errors.New("invalid sessionId")
	}
	session := here.sessions[sessionId]
	count, err := session.mapping.Read(bs)
	if err != nil {
		return 0, err
	}
	return count, nil
}
