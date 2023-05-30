package core

import (
	"ShareMemTCP/util"
	"errors"
)

type ServerSharer struct {
	sessions  map[string]*Session
	selection int8
}

// 创建一个服务端分享者
func NewServerSharer() *ServerSharer {
	return &ServerSharer{
		sessions:  make(map[string]*Session),
		selection: util.SERVER,
	}
}

func (here *ServerSharer) SetCallback() {

}

// 关闭会话
func (here *ServerSharer) Close(sessionId string) error {
	if _, has := here.sessions[sessionId]; !has {
		return errors.New("sessionId error")
	}
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

// 握手函数
func (here *ServerSharer) MemShareTcpDeal(buf []byte, fd int) {

	// if buf[0] == util.MAGIC_NUMBER && buf[1] == util.VERSION {
	// 	s := "yes"
	// 	bs := []byte(s)
	// 	//初始化共享内存

	// 	//init_mem()

	// 	n, err := unix.Write(fd, bs)
	// 	if err != nil {
	// 		fmt.Printf("tcp shake two fail with n = %d\n", n)
	// 	}
	// }
}
