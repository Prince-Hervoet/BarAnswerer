package connection

import (
	"ShareMemTCP/protocol"
	"ShareMemTCP/util"
	"errors"

	"fmt"
	"net"
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
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

// 监听线程tcp连接线程
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

		here.epoll.AddEvent(int32(fd.Fd()), here.memShareTcpDeal)
	}
}

// 握手函数
func (here *Pioneer) memShareTcpDeal(buf []byte, fd int) {

	if buf[0] == util.MAGIC_NUMBER && buf[1] == util.VERSION {
		s := util.SHACK_RESPONSE
		bs := []byte(s)
		//初始化共享内存

		//init_mem()

		n, err := unix.Write(fd, bs)
		if err != nil {
			fmt.Printf("tcp shake two fail with n = %d\n", n)
		}
	}
}

// 打开连接
func (here *Pioneer) OpenConnection(port string, size int32) (string, error) {
	//tcp连接建立
	conn, err := net.Dial("tcp", port)
	if err != nil {
		fmt.Println("error connecting")
		return "", err
	}
	id, _ := gonanoid.New()
	nc := &Connection{
		OpenTimestamp: time.Now().UnixMilli(),
		ConnectionId:  id,
		Conn:          conn,
	}

	here.SheckHand(port, size, conn)

	i := 0
	for ; i < len(here.connections); i++ {
		if here.connections[i] == nil {
			here.connections[i] = nc
			break
		}
	}
	if i >= len(here.connections) {
		here.connections = append(here.connections, nc)
	}

	return id, nil
}

// tcp连接之后的协议栈握手
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
	str := string(readBuf[0:3])

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
	here.connections[which] = nil
}
