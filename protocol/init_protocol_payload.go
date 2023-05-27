package protocol

import "ShareMemTCP/util"

type InitProtocolPayload struct {
	magicNumber int8
	version     int8
	needSize    int32
	pid         int64
}

func (here *InitProtocolPayload) SetMageicNumber(value int8) *InitProtocolPayload {
	here.magicNumber = value
	return here
}

func (here *InitProtocolPayload) SetVersion(value int8) *InitProtocolPayload {
	here.version = value
	return here
}

func (here *InitProtocolPayload) SetNeedSize(value int32) *InitProtocolPayload {
	here.needSize = value
	return here
}

func (here *InitProtocolPayload) SetPid(value int64) *InitProtocolPayload {
	here.pid = value
	return here
}

func (here *InitProtocolPayload) ToByteArray() []byte {
	ans := make([]byte, 14)
	ans[0] = byte(here.magicNumber)
	ans[1] = byte(here.version)
	bs := util.Int32ToBinaryArray(here.needSize)

	for i := 2; i < 6; i++ {
		ans[i] = bs[i-2]
	}

	bs = util.Int64ToBinaryArray(here.pid)
	for i := 6; i < 14; i++ {
		ans[i] = bs[i-6]
	}

	return ans
}
