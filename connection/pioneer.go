package connection

import (
	"ShareMemTCP/core"
	"ShareMemTCP/util"
	"strconv"

	"fmt"
	"net"
)

type Pioneer struct {
	ability byte
	epoll   *EpollInfo
	client  *core.ClientSharer
	server  *core.ServerSharer
}

func (here *Pioneer) ConnectInit(selection byte, port int) (bool, error) {
	here.ability = selection

	//初始化会话结构体
	if selection != util.SERVER {
		here.client = core.NewClientSharer()
	}

	//如果是需要接收数据的话
	if selection > 0 {
		here.server = core.NewServerSharer()

		here.epoll = &EpollInfo{}
		here.epoll.creatEpoll()

		//开启监听线程与epoll处理线程
		go here.Listen(port)
		go here.epollThread()
	}

	return true, nil
}

// 监听线程tcp连接线程
func (here *Pioneer) Listen(port int) {
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

		here.epoll.AddEvent(int32(fd.Fd()), here.server.MemShareTcpDeal)
	}
}

// 打开连接
func (here *Pioneer) OpenConnection(port int, size int32, key int64) {
}

// tcp连接之后的协议栈握手
