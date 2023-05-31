package core

import (
	"ShareMemTCP/util"
	"errors"
)

/*
//dived 分隔符值为'\n'

	type ClinetMessageShakeHand struct {
		Header MessageHeader
		cap    int
	}

	type ServerMessageShakeHand struct {
		Header    MessageHeader
		sessionId string
		dived0     byte
		filePath  string
		dived1     byte
	}

	type ClinetMessageWaveHand struct {
		Header MessageHeader
		sessionId string
		dived0     byte
		filePath  string
		dived1     byte
	}

	type ServerMessageWaveHand struct {
		Header    MessageHeader
		sessionId string
		dived0     byte
		filePath  string
		dived1     byte
	}

	type ClinetMessageNotice struct {
		Header MessageHeader
	}

	type ServerMessageNotice struct {
		Header    MessageHeader
	}
*/
type MessageHeader struct {
	MagicNumber int8
	Version     int8
	CletOrSev   int8
	MessageType int8
	PayLoadLen  int16
}

func CreateMessage(clientOrServer int8, MessageType int8, payLode []byte) []byte {
	offset := 0
	ans := make([]byte, 6)
	ans[offset] = util.MAGIC_NUMBER
	offset++
	ans[offset] = util.VERSION
	offset++
	ans[offset] = byte(clientOrServer)
	offset++
	ans[offset] = byte(MessageType)
	offset++
	if payLode != nil {
		tmp := util.Int16ToBytes(int16(len(payLode)))
		ans[offset] = tmp[0]
		offset++
		ans[offset] = tmp[1]
		ans = append(ans, payLode...)
	}
	return ans
}

func ReadMessege(message []byte, targetMgeType int8, targetSide int8) ([]byte, int16, error) {
	if message[0] != util.MAGIC_NUMBER {
		return nil, 0, errors.New("param is not a messege")
	}
	if message[1] != util.VERSION {
		return nil, 0, errors.New("param is not a messege")
	}
	if message[2] != byte(targetSide) {
		return nil, 0, errors.New("not target messege")
	}
	if message[3] != byte(targetMgeType) {
		return nil, 0, errors.New("not target messege")
	}
	paylen := message[4:6]
	size := util.BytesToInt16(paylen)
	data := message[6 : 6+size]
	return data, size, nil
}

func ReadMessegeHeader(message []byte) (int8, int8, error) {
	if message[0] != util.MAGIC_NUMBER {
		return 0, 0, errors.New("param is not a messege")
	}
	if message[1] != util.VERSION {
		return 0, 0, errors.New("param is not a messege")
	}
	var side int8
	var MessageType int8

	side = int8(message[2])
	MessageType = int8(message[3])
	return side, MessageType, nil
}
