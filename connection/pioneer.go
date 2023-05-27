package connection

import (
	"ShareMemTCP/util"
	"errors"

	//"errors"
	"fmt"
	"net"
	"time"

	"golang.org/x/sys/unix"
)

const (
	CLIENT_ONLY  = 0
	SERVER_ONLY  = 1
	DUPLEX       = 2
	_EVENTS_SIZE = 1024
)

type Pioneer struct {
	connections []*Connection
	id          byte
	epoll       *EpollInfo
	listenFd	int
}

func test() {
	fmt.Println("ce")
}

func (here *Pioneer) NetInit(selection byte) (bool, error) {
	here.id = selection
	//如果是需要接收数据的话
	if selection > 0 {
		here.epoll = &EpollInfo{}
		here.epoll.creatEpoll()
		here.Listen(":20000")

		//初始化map和event数组
		here.epoll.mp = make(map[int32]func())
		here.epoll.events = make([]unix.EpollEvent, _EVENTS_SIZE)
		here.epoll.AddEvent(int32(here.listenFd), test)

		go here.epollThread()
	}

	return true, nil
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

func (here *Pioneer) Listen(port string) (bool, error) {
	// 监听端口

	ln, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println(err)
		return false, errors.New("listen creat fail")
	}
	defer ln.Close()

	file, err := ln.(*net.TCPListener).File()
	if err != nil {
		fmt.Println(err)
		return false, errors.New("listen creat fail")
	}
	here.listenFd = int(file.Fd())
	fmt.Printf("TCP server is listening on port %s, listen fd:%d\n", port, file.Fd())
	return true, nil
}



// 检查连接
func (here *Pioneer) CheckConnection() {

}

// 关闭连接
func (here *Pioneer) CloseConnection(which int) {
	here.connections[which].Conn.Close()
}

// 创造epoll并保存到结构体中
func (here *EpollInfo) creatEpoll() {
	var err error
	here.EpollFd, err = unix.EpollCreate1(unix.EPOLL_CLOEXEC)
	if err != nil {
		return
	}
}

// epoll添加新事件
func (here *EpollInfo) AddEvent(fd int32, function func()) (bool, error) {
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
			unix.Read(int(here.epoll.events[i].Fd), buf)
			(here.epoll.mp[here.epoll.events[i].Fd])() //调用绑定的函数
		}
	}

}
