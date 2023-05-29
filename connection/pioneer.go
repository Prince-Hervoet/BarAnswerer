package connection

import (
	"ShareMemTCP/protocol"
	"ShareMemTCP/util"
	"errors"

	"fmt"
	"net"
	"time"

	"golang.org/x/sys/unix"
)

type Pioneer struct {
	connections []*Connection
	id          byte
	epoll       *EpollInfo
}

func (here *Pioneer) ConnectInit(selection byte, port string) (bool, error) {
	here.id = selection
	//如果是需要接收数据的话
	if selection > 0 {
		here.epoll = &EpollInfo{}
		here.epoll.creatEpoll()

		go here.Listen(port)
		go here.epollThread()
	}

	return true, nil
}

// 创造epoll并保存到结构体中
func (here *EpollInfo) creatEpoll() {
	var err error
	here.EpollFd, err = unix.EpollCreate1(unix.EPOLL_CLOEXEC)
	if err != nil {
		return
	}
	//初始化map和event数组
	here.mp = make(map[int32]func([]byte, int))
	here.events = make([]unix.EpollEvent, util.EVENTS_SIZE)
}

// epoll添加新事件
func (here *EpollInfo) AddEvent(fd int32, function func([]byte, int)) (bool, error) {
	// 添加文件描述符到 epoll 实例并监听可读事件
	var event unix.EpollEvent
	event.Events = unix.EPOLLIN
	event.Fd = fd // 文件描述符
	if err := unix.EpollCtl(here.EpollFd, unix.EPOLL_CTL_ADD, int(fd), &event); err != nil {
		fmt.Printf("Error adding file descriptor to epoll instance: %v\n", err)
		unix.Close(here.EpollFd)
		return false, errors.New("add epoll event fail")
	}
	here.mp[fd] = function
	return true, nil
}

func (here *EpollInfo) DeleteEvent(fd int32) (bool, error) {

	if err := unix.EpollCtl(here.EpollFd, unix.EPOLL_CTL_DEL, int(fd), nil); err != nil {
		fmt.Printf("Error delete file descriptor to epoll instance: %v\n", err)
		unix.Close(here.EpollFd)
		return false, errors.New("delete epoll event fail")
	}

	delete(here.mp, fd)
	return true, nil
}

func (here *Pioneer) memShareTcpDeal(buf []byte, fd int) {
	//fmt.Println(buf)
	if buf[0] == util.MAGIC_NUMBER && buf[1] == util.VERSION {
		//fmt.Println("rev frame")
		s := "yes"
		bs := []byte(s)
		//初始化共享内存

		//init_mem()

		fmt.Printf("deal fd:%d\n", fd)
		unix.Write(fd, bs)
	}
}

func (here *Pioneer) Listen(port string) {
	// 监听端口
	ln, err := net.Listen("tcp", port)
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

	fmt.Printf("TCP server is listening on port %s, listen fd:%d\n", port, file.Fd())

	//buf := make([]byte, 4)
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
		//unix.Read(int(fd.Fd()), buf)
		//fmt.Println(buf)
		here.epoll.AddEvent(int32(fd.Fd()), here.memShareTcpDeal)
	}
}

func (here *Pioneer) epollThread() {
	buf := make([]byte, 1024)
	for {
		// 等待事件发生
		n, err := unix.EpollWait(here.epoll.EpollFd, here.epoll.events, -1)
		if err != nil {
			fmt.Printf("Error waiting for events: %v\n", err)
			unix.Close(here.epoll.EpollFd)
			return
		}
		defer unix.Close(here.epoll.EpollFd)

		// 处理事件
		for i := 0; i < n; i++ {
			fd := here.epoll.events[i].Fd
			unix.Read(int(fd), buf)
			(here.epoll.mp[fd])(buf, int(fd)) //调用绑定的函数
		}
	}
}

// 打开连接
func (here *Pioneer) OpenConnection(port string, size int32) (int64, error) {
	//tcp连接建立
	conn, err := net.Dial("tcp", port)
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

	here.SheckHand(port, size, conn)

	here.connections = append(here.connections, nc)
	return id, nil
}

func (here *Pioneer) SheckHand(port string, size int32, conn net.Conn) (bool, error) {
	//共享内存协议栈
	initpro := protocol.InitProtocolPayload{}
	initpro.SetMageicNumber(util.MAGIC_NUMBER)
	initpro.SetNeedSize(size)
	initpro.SetVersion(util.VERSION)

	//发送协议头
	buf := initpro.ToByteArray()
	fmt.Println(buf)
	n, err := conn.Write(buf)
	if err != nil || n < 0 {
		fmt.Println("error menTCP connecting")
		return false, err
	}

	//等待服务端初始化完成
	readBuf := make([]byte, 32)
	conn.Read(readBuf)
	str := string(readBuf)
	if str != util.SHACK_RESPONSE {
		fmt.Println("memTcp response error")
		return false, errors.New("memTcp response error")
	}
	fmt.Println("memTcp succes")
	return true, nil
}

// 检查连接
func (here *Pioneer) CheckConnection() {

}

// 关闭连接
func (here *Pioneer) CloseConnection(which int) {

	//断开tcp连接
	here.connections[which].Conn.Close()
}
