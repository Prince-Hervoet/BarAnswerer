package memory

import "ShareMemTCP/util"

type ShareMemoryHeader struct {
	status int8
	head   int32
	tail   int32
	size   int32
	cap    int32
}

func (here *ShareMemoryHeader) FromBytes(data []byte) {
	if len(data) < SHARE_MEMORY_HEADER_SIZE {
		return
	}
	here.status = int8(data[0])
	here.head = util.BytesToInt32(data[1:5])
	here.tail = util.BytesToInt32(data[5:9])
	here.size = util.BytesToInt32(data[9:13])
	here.cap = util.BytesToInt32(data[13:17])
}

func (here *ShareMemoryHeader) ToBytes() []byte {
	ans := make([]byte, 0, 17)
	ans = append(ans, byte(here.status))
	bs := util.Int32ToBytes(here.head)
	ans = append(ans, bs...)
	bs = util.Int32ToBytes(here.tail)
	ans = append(ans, bs...)
	bs = util.Int32ToBytes(here.size)
	ans = append(ans, bs...)
	bs = util.Int32ToBytes(here.cap)
	ans = append(ans, bs...)
	return ans
}
