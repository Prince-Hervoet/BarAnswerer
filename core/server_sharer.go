package core

import (
	// "ShareMemTCP/memory"
	// "encoding/json"
	"ShareMemTCP/util"
	"errors"
	// "fmt"
	// "net"
	// gonanoid "github.com/matoous/go-nanoid/v2"
)

type ServerSharer struct {
	sessions  map[int64]*Session
	selection int8
}

// 创建一个服务端分享者
func NewServerSharer() *ServerSharer {
	return &ServerSharer{
		sessions:  make(map[int64]*Session),
		selection: util.SERVER,
	}
}

func (here *ServerSharer) SetCallback() {

}

// 关闭会话
func (here *ServerSharer) Close(key int64) error {
	if _, has := here.sessions[key]; !has {
		return errors.New("key error")
	}
	return nil
}

func (here *ServerSharer) Reply(content []byte, key int64) error {
	if _, has := here.sessions[key]; !has {
		return errors.New("sessionId error")
	}
	session := here.sessions[key]
	session.mapping.Write(content)
	return nil
}

func (here *ServerSharer) Read() {

}

// 开启监听端口，等待握手事件到来然后开启共享内存
// func (here *ServerSharer) Open(address *Address) {
// 	localAddr := net.UDPAddr{
// 		IP:   net.ParseIP("127.0.0.1"),
// 		Port: address.Port,
// 	}
// 	conn, err := net.ListenUDP("udp", &localAddr)
// 	if err != nil {
// 		fmt.Printf("Error creating UDP listener: %v\n", err)
// 		return
// 	}
// 	buffer := make([]byte, 256)
// 	fmt.Println("UDP server running")
// 	defer func() {
// 		conn.Close()
// 	}()
// 	for {
// 		_, addr, err := conn.ReadFromUDP(buffer)
// 		if err != nil {
// 			fmt.Printf("Error reading UDP packet: %v\n", err)
// 			continue
// 		}

// 		message, err := ReadJsonByLineFeed(buffer)
// 		if err != nil || !CheckHeader(message) {
// 			if err != nil {
// 				fmt.Println(err.Error())
// 			}
// 			continue
// 		}

// 		switch message.Header.MessageType {
// 		case MESSAGE_TYPE_INIT:
// 			fb := message.Payload.(*FirstBody)
// 			sm := memory.OpenShareMemory()
// 			id, _ := gonanoid.New()
// 			filePath, err := sm.OpenFile(id, fb.NeedCap)
// 			if err != nil {
// 				fmt.Println(err.Error())
// 				sm.Close()
// 				continue
// 			}
// 			sb := &SecondBody{
// 				Header: &MessageHeader{
// 					MagicNumber: MAGIC_NUMBER,
// 					Version:     VERSION,
// 					MessageType: MESSAGE_TYPE_INIT,
// 				},
// 				FilePath: filePath,
// 				Cap:      fb.NeedCap,
// 			}

// 			bs, err := json.Marshal(sb)
// 			if err != nil {
// 				fmt.Println(err.Error())
// 				sm.Close()
// 				continue
// 			}
// 			_, err = conn.WriteToUDP(bs, addr)
// 			if err != nil {
// 				fmt.Printf("Error sending response over UDP: %v\n", err)
// 				continue
// 			}
// 			here.sessions[fb.Key] = NewSession(fb.Key, address, sm, conn)
// 		case MESSAGE_TYPE_DATA:
// 		case MESSAGE_TYPE_QUIT:
// 		}

// 	}
// }
