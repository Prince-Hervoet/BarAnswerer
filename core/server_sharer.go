package core

import (
	"ShareMemTCP/util"
	"errors"
)

type ServerSharer struct {
	sessions  map[int64]*Session
	selection int8
}

// 创建一个服务端分享者
func NewServerSharer() *ServerSharer {
	return &ServerSharer{
		sessions:  make(map[int64]*Session),
		selection: util.SERVER,
	}
}

func (here *ServerSharer) SetCallback() {

}

// 关闭会话
func (here *ServerSharer) Close(key int64) error {
	if _, has := here.sessions[key]; !has {
		return errors.New("key error")
	}
	return nil
}

// 当客户端发来请求时的回复函数
func (here *ServerSharer) Reply(content []byte, key int64) error {
	if _, has := here.sessions[key]; !has {
		return errors.New("sessionId error")
	}
	session := here.sessions[key]
	session.mapping.Write(content)
	return nil
}

func (here *ServerSharer) Read() {

}

// 开启监听端口，等待握手事件到来然后开启共享内存
func (here *ServerSharer) Open(address *Address) {

}


