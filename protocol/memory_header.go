package protocol

import "ShareMemTCP/util"

type MemoryHeader struct {
	Flag int8
	Head int32
	Tail int32
}

func NewMemoryHeader(flag int8, head, tail int32) *MemoryHeader {
	return &MemoryHeader{
		Flag: flag,
		Head: head,
		Tail: tail,
	}
}

func (here *MemoryHeader) FromByteArray(data []byte) {
	if len(data) < 9 {
		return
	}
	here.Flag = int8(data[0])
	bs := data[1:5]
	here.Head = util.BytesToInt32(bs)
	bs = data[5:9]
	here.Tail = util.BytesToInt32(bs)
}

func (here *MemoryHeader) ToByteArray() []byte {
	ans := make([]byte, 9)
	ans[0] = byte(here.Flag)
	bs := util.Int32ToBytes(here.Head)
	ans[1] = bs[0]
	ans[2] = bs[1]
	ans[3] = bs[2]
	ans[4] = bs[3]

	bs = util.Int32ToBytes(here.Tail)
	ans[5] = bs[0]
	ans[6] = bs[1]
	ans[7] = bs[2]
	ans[8] = bs[3]
	return ans
}
