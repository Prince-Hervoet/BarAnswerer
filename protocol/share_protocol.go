package protocol

import "ShareMemTCP/util"

type ShareProtocol struct {
	MagicNumber int8
	MessageType int8
	PayloadSize int32
	Payload     []byte
}

func NewShareProtocol() *ShareProtocol {
	return &ShareProtocol{
		MagicNumber: util.MAGIC_NUMBER,
	}
}

func (here *ShareProtocol) ToByteArray() []byte {
	ans := make([]byte, 0)
	b1 := byte(here.MagicNumber)
	b2 := byte(here.MessageType)
	b3 := util.Int32ToBytes(here.PayloadSize)
	ans = append(ans, b1, b2)
	for i := 0; i < 4; i++ {
		ans = append(ans, b3[i])
	}
	ans = append(ans, here.Payload...)
	return ans
}

func (here *ShareProtocol) FromByteArray(data []byte) {
	if len(data) < 14 {
		return
	}
	v1 := data[0]
	v2 := data[1]
	here.MagicNumber = int8(v1)
	here.MessageType = int8(v2)

	v3 := data[2:6]
	here.PayloadSize = util.BytesToInt32(v3)

	here.Payload = data[6:]
}

func (here *ShareProtocol) SetMageicNumber(value int8) *ShareProtocol {
	here.MessageType = value
	return here
}

func (here *ShareProtocol) SetMessageType(value int8) *ShareProtocol {
	here.MessageType = value
	return here
}

func (here *ShareProtocol) SetPayloadSize(value int32) *ShareProtocol {
	here.PayloadSize = value
	return here
}

func (here *ShareProtocol) SetPayload(value []byte) *ShareProtocol {
	here.Payload = value
	return here
}
