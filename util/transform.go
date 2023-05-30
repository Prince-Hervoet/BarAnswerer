package util

import "encoding/binary"

func BytesToInt32(data []byte) int32 {
	return int32(binary.BigEndian.Uint32(data))
}

func Int32ToBytes(num int32) []byte {
	ans := make([]byte, 4)
	binary.BigEndian.PutUint32(ans, uint32(num))
	return ans
}

func Int32Min(a, b int32) int32 {
	if a > b {
		return b
	}
	return a
}

func IntMin(a, b int) int {
	if a > b {
		return b
	}
	return a
}
