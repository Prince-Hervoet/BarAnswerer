package core

type Pioneer struct {
	Client *ClientSharer
	Server *ServerSharer
	epoll  *EpollInfo
}

func (here *Pioneer) ConnectInit(port int) (bool, error) {
	//初始化会话结构体
	here.epoll = CreatEpoll()
	here.Client = NewClientSharer()
	here.Server = NewServerSharer()
	here.Client.EpollIndex = here.epoll
	here.Server.EpollIndex = here.epoll

	//开启监听线程与epoll处理线程
	go here.Server.Listen(port)
	go here.epollThread()

	//如果是需要接收数据的话

	return true, nil
}

func (here *Pioneer) NewPoineer(port int) *Pioneer {
	t := &Pioneer{}
	t.ConnectInit(port)
	return t
}

// 设置服务端读取客户端的数据后的回调函数
func (here *Pioneer) SetCallback(sessionId string, call func([]byte)) {
	here.Server.callBacks[sessionId] = call
}

// 打开连接
func (here *Pioneer) OpenConnection(port int, size int32, key int64) {
}

// tcp连接之后的协议栈握手
