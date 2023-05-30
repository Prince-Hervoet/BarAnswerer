package core

import (
	// "ShareMemTCP/memory"
	//"ShareMemTCP/util"
	// "encoding/json"
	"ShareMemTCP/util"
	"errors"
	// "fmt"
	// "net"
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
	// isLinked := true
	// remoteAddr := &net.UDPAddr{
	// 	IP:   net.ParseIP("127.0.0.1"),
	// 	Port: address.Port,
	// }
	// conn, err := net.DialUDP("udp", nil, remoteAddr)
	// if err != nil {
	// 	fmt.Printf("Error creating UDP connection: %v\n", err)
	// 	return
	// }
	// defer func() {
	// 	if !isLinked {
	// 		conn.Close()
	// 	}
	// }()

	// fb := &FirstBody{
	// 	Header: &MessageHeader{
	// 		MagicNumber: util.MAGIC_NUMBER,
	// 		Version:     1,
	// 		MessageType: util.MESSAGE_TYPE_INIT,
	// 	},
	// 	NeedCap: cap,
	// 	Key:     key,
	// }

	// bs, _ := json.Marshal(fb)
	// bs = append(bs, '\n')
	// _, err = conn.Write(bs)
	// if err != nil {
	// 	isLinked = false
	// 	return
	// }
	// buffer := make([]byte, 256)
	// count, err := conn.Read(buffer)
	// if err != nil {
	// 	isLinked = false
	// 	fmt.Printf("Error sending data over UDP: %v\n", err)
	// 	return
	// }

	// sb := &SecondBody{}
	// err = json.Unmarshal(buffer[0:count], sb)
	// if err != nil {
	// 	isLinked = false
	// 	fmt.Println(err.Error())
	// 	return
	// }

	// sm := memory.OpenShareMemory()
	// err = sm.LinkFile(sb.FilePath, sb.Cap)
	// if err != nil {
	// 	isLinked = false
	// 	sm.Close()
	// 	return
	// }
	// fmt.Println(sb)
	// here.sessions[key] = NewSession(key, address, sm, conn)
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

// 　读取共享内存中数据
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

// func (here *ClientSharer) UDPSend(remotePort *Address, key int64, data []byte) {
// 	remoteAddr := net.UDPAddr{
// 		IP:   net.ParseIP("127.0.0.1"),
// 		Port: remotePort.Port,
// 	}
// 	conn, err := net.DialUDP("udp", nil, &remoteAddr)
// 	if err != nil {
// 		fmt.Printf("Error creating UDP connection: %v\n", err)
// 		return
// 	}
// 	defer conn.Close()

// 	_, _ = conn.Write(data)
// 	buffer := make([]byte, 256)
// 	count, err := conn.Read(buffer)
// 	if err != nil {
// 		fmt.Printf("Error sending data over UDP: %v\n", err)
// 		return
// 	}

// 	sb := &SecondBody{}
// 	err = json.Unmarshal(buffer[0:count], sb)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		return
// 	}

// 	sm := memory.OpenShareMemory()
// 	err = sm.LinkFile(sb.FilePath, sb.Cap)
// 	if err != nil {
// 		return
// 	}
// 	here.sessions[key] = NewSession(key, remotePort, sm, nil)
// }
