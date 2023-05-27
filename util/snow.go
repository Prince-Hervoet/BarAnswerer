package util

import (
	"fmt"
	"time"
)

const (
	// 雪花算法的起始时间戳，可以根据实际情况进行调整
	twepoch = int64(1586921080000)

	// 机器ID所占的位数
	workerBits = uint(5)

	// 数据中心ID所占的位数
	datacenterBits = uint(5)

	// 序列号所占的位数
	sequenceBits = uint(12)

	// 每一部分的最大值
	maxWorkerId     = -1 ^ (-1 << workerBits)
	maxDatacenterId = -1 ^ (-1 << datacenterBits)
	maxSequence     = -1 ^ (-1 << sequenceBits)

	// 每一部分向左的位移
	workerShift        = sequenceBits
	datacenterShift    = sequenceBits + workerBits
	timestampLeftShift = sequenceBits + workerBits + datacenterBits

	// 工作机器ID
	workerId int64 = 1

	// 数据中心ID
	datacenterId int64 = 1
)

var lastTimestamp int64
var sequence int64

// 生成一个新的Snowflake ID
func NextId() int64 {
	// 获取当前时间戳
	currentTime := time.Now().UnixNano() / 1e6

	if currentTime < lastTimestamp {
		panic(fmt.Sprintf("clock moved backwards, refusing to generate id for %d milliseconds", lastTimestamp-currentTime))
	}

	if currentTime == lastTimestamp {
		sequence = (sequence + 1) & maxSequence

		if sequence == 0 {
			for currentTime <= lastTimestamp {
				currentTime = time.Now().UnixNano() / 1e6
			}
		}
	} else {
		sequence = 0
	}

	lastTimestamp = currentTime

	return ((currentTime - twepoch) << timestampLeftShift) |
		(datacenterId << datacenterShift) |
		(workerId << workerShift) |
		sequence
}
