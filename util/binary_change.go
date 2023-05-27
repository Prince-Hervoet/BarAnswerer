package util

import "encoding/binary"

func Int32ToBinaryArray(number int32) []byte {
	ans := make([]byte, 4)
	binary.BigEndian.PutUint32(ans, uint32(number))
	return ans
}

func Int64ToBinaryArray(number int64) []byte {
	ans := make([]byte, 8)
	binary.BigEndian.PutUint64(ans, uint64(number))
	return ans
}
