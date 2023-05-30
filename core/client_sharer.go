package core

import (
	"ShareMemTCP/memory"
	"encoding/json"
	"errors"
	"fmt"
	"net"
)

type ClientSharer struct {
	sessions  map[int64]*Session
	selection int8
}

func NewClientSharer() *ClientSharer {
	return &ClientSharer{
		sessions:  make(map[int64]*Session),
		selection: CLIENT,
	}
}

func (here *ClientSharer) Test() {
	fb := &FirstBody{
		Header: &MessageHeader{
			MagicNumber: MAGIC_NUMBER,
			Version:     1,
			MessageType: MESSAGE_TYPE_INIT,
		},
		NeedCap: 4096,
	}
	b, _ := json.Marshal(fb)
	b = append(b, '\n')
	go here.UDPSend(&Address{
		Ip:   "",
		Port: 5551,
	}, 13, b)
}

func (here *ClientSharer) Link(address *Address, cap int32, key int64) {
	isLinked := true
	remoteAddr := &net.UDPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: address.Port,
	}
	conn, err := net.DialUDP("udp", nil, remoteAddr)
	if err != nil {
		fmt.Printf("Error creating UDP connection: %v\n", err)
		return
	}
	defer func() {
		if !isLinked {
			conn.Close()
		}
	}()

	fb := &FirstBody{
		Header: &MessageHeader{
			MagicNumber: MAGIC_NUMBER,
			Version:     1,
			MessageType: MESSAGE_TYPE_INIT,
		},
		NeedCap: cap,
		Key:     key,
	}

	bs, _ := json.Marshal(fb)
	bs = append(bs, '\n')
	_, err = conn.Write(bs)
	if err != nil {
		isLinked = false
		return
	}
	buffer := make([]byte, 256)
	count, err := conn.Read(buffer)
	if err != nil {
		isLinked = false
		fmt.Printf("Error sending data over UDP: %v\n", err)
		return
	}

	sb := &SecondBody{}
	err = json.Unmarshal(buffer[0:count], sb)
	if err != nil {
		isLinked = false
		fmt.Println(err.Error())
		return
	}

	sm := memory.OpenShareMemory()
	err = sm.LinkFile(sb.FilePath, sb.Cap)
	if err != nil {
		isLinked = false
		sm.Close()
		return
	}
	fmt.Println(sb)
	here.sessions[key] = NewSession(key, address, sm, conn)
}

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

func (here *ClientSharer) AddAcceptor(handler func(message []byte)) {

}

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

func (here *ClientSharer) SyncRead(bs []byte, sessionId int64) (int32, error) {
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

func (here *ClientSharer) UDPSend(remotePort *Address, key int64, data []byte) {
	remoteAddr := net.UDPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: remotePort.Port,
	}
	conn, err := net.DialUDP("udp", nil, &remoteAddr)
	if err != nil {
		fmt.Printf("Error creating UDP connection: %v\n", err)
		return
	}
	defer conn.Close()

	_, err = conn.Write(data)
	buffer := make([]byte, 256)
	count, err := conn.Read(buffer)
	if err != nil {
		fmt.Printf("Error sending data over UDP: %v\n", err)
		return
	}

	sb := &SecondBody{}
	err = json.Unmarshal(buffer[0:count], sb)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	sm := memory.OpenShareMemory()
	err = sm.LinkFile(sb.FilePath, sb.Cap)
	if err != nil {
		return
	}
	here.sessions[key] = NewSession(key, remotePort, sm, nil)
}
