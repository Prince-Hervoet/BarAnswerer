package core

import (
	"ShareMemTCP/util"
	"errors"
	"fmt"

	"golang.org/x/sys/unix"
)

type EpollInfo struct {
	EpollFd  int
	EventNum int
	events   []unix.EpollEvent
	mp       map[int32]func([]byte, int)
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
	// 添加文件描述符到 epoll 实例并监听可读或者tcp断开事件
	var event unix.EpollEvent
	event.Events = unix.EPOLLIN | unix.EPOLLHUP
	event.Fd = fd // 文件描述符

	if err := unix.EpollCtl(here.EpollFd, unix.EPOLL_CTL_ADD, int(fd), &event); err != nil {
		fmt.Printf("Error adding file descriptor to epoll instance: %v\n", err)
		unix.Close(here.EpollFd)
		return false, errors.New("add epoll event fail")
	}

	//映射fd与回调函数
	here.mp[fd] = function
	return true, nil
}

// 删除epoll中监管fd
func (here *EpollInfo) DeleteEvent(fd int32) (bool, error) {

	if err := unix.EpollCtl(here.EpollFd, unix.EPOLL_CTL_DEL, int(fd), nil); err != nil {
		fmt.Printf("Error delete file descriptor to epoll instance: %v\n", err)
		unix.Close(here.EpollFd)
		return false, errors.New("delete epoll event fail")
	}

	delete(here.mp, fd)
	return true, nil
}

// epoll处理线程
func (here *ServerSharer) epollThread() {
	for {
		buf := make([]byte, 1024)
		// 等待事件发生
		n, err := unix.EpollWait(here.Epoll.EpollFd, here.Epoll.events, -1)
		if err != nil {
			fmt.Printf("Error waiting for events: %v\n", err)
			unix.Close(here.Epoll.EpollFd)
			return
		}
		defer unix.Close(here.Epoll.EpollFd)

		// 处理事件
		for i := 0; i < n; i++ {
			fd := here.Epoll.events[i].Fd
			n, err := unix.Read(int(fd), buf)
			if n == 0 || err != nil {
				fmt.Printf("link fd %d is over now delete it\n", fd)
				//删除对方的client的相关对话资源
				here.Epoll.DeleteEvent(fd)                          //移除epoll监控
				here.sessions[here.SidMap[int(fd)]].mapping.Close() //移除共享内存映射
				here.DeleteSession(int(fd))                         //移除session映射
				unix.Close(int(fd))

				//删除本进程的client的相关对话资源
				

				continue
			}
			go (here.Epoll.mp[fd])(buf, int(fd)) //调用绑定的函数
		}
	}
}
