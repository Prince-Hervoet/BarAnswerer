package core

type Pioneer struct {
	Client  *ClientSharer
	Server  *ServerSharer
}

func (here *Pioneer) ConnectInit(port int) (bool, error) {
	//初始化会话结构体
	here.Client = NewClientSharer()
	here.Server = NewServerSharer()
	here.Server.Epoll.creatEpoll()

	//开启监听线程与epoll处理线程
	go here.Server.Listen(port)
	go here.Server.epollThread()

	//如果是需要接收数据的话

	return true, nil
}

func (here *Pioneer) NewPoineer(port int)(*Pioneer){
	t := &Pioneer{}
	t.ConnectInit(port)
	return t
}

// 打开连接
func (here *Pioneer) OpenConnection(port int, size int32, key int64) {
}

// tcp连接之后的协议栈握手
