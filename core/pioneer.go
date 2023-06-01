package core

type Pioneer struct {
	client *ClientSharer
	server *ServerSharer
	epoll  *EpollInfo
}

func NewPoineer(port int) *Pioneer {
	t := &Pioneer{}
	t.connectInit(port)
	return t
}

// 打开连接
func (here *Pioneer) Link(port int, cap int32) (string, error) {
	sid, err := here.client.Link(port, cap)
	if err != nil {
		return "", err
	}
	return sid, nil
}

// 设置服务端读取客户端的数据后的回调函数
func (here *Pioneer) SetCallback(sessionId string, call func([]byte)) {
	here.server.callBacks[sessionId] = call
}

// 写数据
func (here *Pioneer) Write(data []byte, sessionId string) error {
	return here.client.Send(data, sessionId)
}

// 关闭连接
func (here *Pioneer) Close(sessionId string) error {
	return here.client.Close(sessionId)
}

func (here *Pioneer) connectInit(port int) (bool, error) {
	//初始化会话结构体
	here.epoll = CreatEpoll()
	here.client = newClientSharer()
	here.server = newServerSharer()
	here.client.EpollIndex = here.epoll
	here.server.EpollIndex = here.epoll

	//开启监听线程与epoll处理线程
	go here.server.Listen(port)
	go here.epollThread()

	//如果是需要接收数据的话
	return true, nil
}
