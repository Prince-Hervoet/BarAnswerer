package core

import (
	"ShareMemTCP/util"
)

type Pioneer struct {
	ability byte
	client  *ClientSharer
	server  *ServerSharer
}

func (here *Pioneer) ConnectInit(selection byte, port int) (bool, error) {
	here.ability = selection

	//初始化会话结构体
	if selection != util.SERVER {
		here.client = NewClientSharer()
	}

	//如果是需要接收数据的话
	if selection > 0 {
		here.server = NewServerSharer()

		here.server.Epoll.creatEpoll()

		//开启监听线程与epoll处理线程
		go here.server.Listen(port)
		go here.server.epollThread()
	}

	return true, nil
}

// 打开连接
func (here *Pioneer) OpenConnection(port int, size int32, key int64) {
}

// tcp连接之后的协议栈握手
