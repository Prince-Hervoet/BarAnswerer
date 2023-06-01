package core

type Pioneer struct {
	Client *ClientSharer
	Server *ServerSharer
	epoll  *EpollInfo
}

func (here *Pioneer) NewPoineer(port int) *Pioneer {
	t := &Pioneer{}
	t.ConnectInit(port)
	return t
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

// 设置服务端读取客户端的数据后的回调函数
func (here *Pioneer) SetReadCallback(sessionId string, call func([]byte)) {
	here.Server.callBacks[sessionId] = call
}

// 打开连接
func (here *Pioneer) OpenConnection(port int, size int32) (string, error) {
	seesionId, err := here.Client.Link(port, size)
	return seesionId, err
}

// 写数据
func (here *Pioneer) Write(data []byte, sessionId string) error {
	return here.Client.Send(data, sessionId)
}

// 关闭连接
func (here *Pioneer) Close(sessionId string) error {
	return here.Client.Close(sessionId)
}
