package protocol

import "ShareMemTCP/util"

type MemoryHeader struct {
	Flag int8
	Head int32
	Tail int32
	Size int32
	Cap  int32
}

func NewMemoryHeader(flag int8, head, tail, Size, Cap int32) *MemoryHeader {
	return &MemoryHeader{
		Flag: flag,
		Head: head,
		Tail: tail,
		Size: Size,
		Cap:  Cap,
	}
}

func (here *MemoryHeader) FromByteArray(data []byte) {
	if len(data) < util.MEM_HEADER_SIZE {
		return
	}
	here.Flag = int8(data[0])
	bs := data[1:5]
	here.Head = util.BinaryArrayToInt32(bs)
	bs = data[5:9]
	here.Tail = util.BinaryArrayToInt32(bs)
	bs = data[9:13]
	here.Size = util.BinaryArrayToInt32(bs)
	bs = data[13:17]
	here.Cap = util.BinaryArrayToInt32(bs)
}

func (here *MemoryHeader) ToByteArray() []byte {
	ans := make([]byte, 9)
	ans[0] = byte(here.Flag)
	bs := util.Int32ToBinaryArray(here.Head)
	ans[1] = bs[0]
	ans[2] = bs[1]
	ans[3] = bs[2]
	ans[4] = bs[3]

	bs = util.Int32ToBinaryArray(here.Tail)
	ans[5] = bs[0]
	ans[6] = bs[1]
	ans[7] = bs[2]
	ans[8] = bs[3]

	bs = util.Int32ToBinaryArray(here.Size)
	ans[9] = bs[0]
	ans[10] = bs[1]
	ans[11] = bs[2]
	ans[12] = bs[3]

	bs = util.Int32ToBinaryArray(here.Cap)
	ans[13] = bs[0]
	ans[14] = bs[1]
	ans[15] = bs[2]
	ans[16] = bs[3]
	return ans
}
