package util

import "encoding/binary"

func BytesToInt16(data []byte) int16 {
	return int16(binary.BigEndian.Uint16(data))
}

func Int16ToBytes(num int16) []byte {
	ans := make([]byte, 4)
	binary.BigEndian.PutUint16(ans, uint16(num))
	return ans
}

func BytesToInt32(data []byte) int32 {
	return int32(binary.BigEndian.Uint32(data))
}

func Int32ToBytes(num int32) []byte {
	ans := make([]byte, 4)
	binary.BigEndian.PutUint32(ans, uint32(num))
	return ans
}

func Int64ToBytes(number int64) []byte {
	ans := make([]byte, 8)
	binary.BigEndian.PutUint64(ans, uint64(number))
	return ans
}

func BytesToInt64(data []byte) int64 {
	return int64(binary.BigEndian.Uint64(data))
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
