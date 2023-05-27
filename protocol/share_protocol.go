package protocol

import "ShareMemTCP/util"

type ShareProtocol struct {
	magicNumber int8
	messageType int8
	payloadSize int32
	payload     []byte
}

func NewShareProtocol() *ShareProtocol {
	return &ShareProtocol{
		magicNumber: MAGIC_NUMBER,
	}
}

func (here *ShareProtocol) ToByteArray() []byte {
	ans := make([]byte, 0)
	b1 := byte(here.magicNumber)
	b2 := byte(here.messageType)
	b3 := util.Int32ToBinaryArray(here.payloadSize)
	ans = append(ans, b1, b2)
	for i := 0; i < 4; i++ {
		ans = append(ans, b3[i])
	}
	ans = append(ans, here.payload...)
	return ans
}

func (here *ShareProtocol) FromByteArray(data []byte) {
	if len(data) < 14 {
		return
	}
	v1 := data[0]
	v2 := data[1]
	here.magicNumber = int8(v1)
	here.messageType = int8(v2)

	v3 := data[2:6]
	here.payloadSize = util.BinaryArrayToInt32(v3)

	here.payload = data[6:]
}

func (here *ShareProtocol) SetMageicNumber(value int8) *ShareProtocol {
	here.messageType = value
	return here
}

func (here *ShareProtocol) SetMessageType(value int8) *ShareProtocol {
	here.messageType = value
	return here
}

func (here *ShareProtocol) SetPayloadSize(value int32) *ShareProtocol {
	here.payloadSize = value
	return here
}

func (here *ShareProtocol) SetPayload(value []byte) *ShareProtocol {
	here.payload = value
	return here
}
