package memory

import (
	"ShareMemTCP/util"
	"errors"
)

type ActionSection struct {
	head    int
	tail    int
	section []byte
	cap     int
	size    int
}

func NewActionSection(cap int, section []byte) *ActionSection {
	return &ActionSection{
		cap:     cap,
		section: section,
	}
}

func (here *ActionSection) Write(data []byte) error {
	if len(data)+here.size > here.cap {
		return errors.New("no enough error")
	}
	for i := 0; i < len(data); i++ {
		here.section[here.tail] = data[i]
		if here.tail == here.size-1 {
			here.tail = 0
		} else {
			here.tail += 1
		}
	}
	here.size += len(data)
	return nil
}

func (here *ActionSection) Read(len int) []byte {
	if here.size == 0 {
		return make([]byte, 0)
	}
	ansLen := util.IntMin(len, here.size)
	ans := make([]byte, ansLen)
	for i := 0; i < ansLen; i++ {
		ans[i] = here.section[here.head]
		if here.head == here.size-1 {
			here.head = 0
		} else {
			here.head += 1
		}
	}
	here.size -= ansLen
	return ans
}

func (here *ActionSection) Size() int {
	return here.size
}
