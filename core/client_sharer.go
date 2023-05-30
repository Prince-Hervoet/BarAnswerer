package core

import (
	"ShareMemTCP/util"
	"errors"
)

type ClientSharer struct {
	sessions  map[int64]*Session
	selection int8
}

// 创建客户端分享者
func NewClientSharer() *ClientSharer {
	return &ClientSharer{
		sessions:  make(map[int64]*Session),
		selection: util.CLIENT,
	}
}

// 链接对端进程，请求开启共享内存（传入的对端进程的地址、需要的内存大小、连接的标识）
func (here *ClientSharer) Link(address *Address, cap int32, key int64) {
	//tcp连接建立
	// conn, err := net.Dial("tcp", address.Port)
	// if err != nil {
	// 	fmt.Println("error connecting")
	// 	return "", err
	// }
	// id, _ := gonanoid.New()

	// here.SheckHand(port, size, conn)

	// return id, nil
}

// 发送数据到共享内存中（需要通知对端进程进行读取）
func (here *ClientSharer) Send(data []byte, key int64) error {
	if _, has := here.sessions[key]; !has {
		return errors.New("invalid sessionId")
	}
	session := here.sessions[key]
	err := session.mapping.Write(data)
	if err != nil {
		return err
	}
	return nil
}

// 读取共享内存中数据
func (here *ClientSharer) Read(bs []byte, sessionId int64) (int32, error) {
	if _, has := here.sessions[sessionId]; !has {
		return 0, errors.New("invalid sessionId")
	}
	session := here.sessions[sessionId]
	count, err := session.mapping.Read(bs)
	if err != nil {
		return 0, err
	}
	return count, nil
}
