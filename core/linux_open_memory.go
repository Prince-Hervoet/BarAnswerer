package core

import (
	"ShareMemTCP/util"
	"errors"
	"fmt"
	"syscall"
	"unsafe"
)

const (
	IPC_CREATE = 00001000
	PAGE_SIZE  = 4096
)

type ShareMemory struct {
	leftLimit  uintptr
	rightLimit uintptr
	head       uintptr
	tail       uintptr
	size       int
	cap        int
}

func OpenShareMemory(key int, cap int) *ShareMemory {
	shmid, _, err := syscall.Syscall(syscall.SYS_SHMGET, uintptr(key), uintptr(112), IPC_CREATE|0666)
	if err != 0 {
		fmt.Printf("open memory error: %s\n", err.Error())
		return nil
	}
	shmaddr, _, err := syscall.Syscall(syscall.SYS_SHMAT, shmid, 0, 0)
	if err != 0 {
		fmt.Printf("syscall error, err: %v\n", err)
		return nil
	}
	fmt.Println(shmaddr)
	return &ShareMemory{
		leftLimit:  shmaddr,
		rightLimit: shmaddr + uintptr(cap),
		head:       shmaddr,
		tail:       shmaddr,
		size:       0,
		cap:        cap,
	}
}

func (here *ShareMemory) WriteShareMemory(data []byte) error {
	if len(data)+here.size > here.cap {
		return errors.New("memory no enough")
	}
	for i := 0; i < len(data); i++ {
		*(*byte)(unsafe.Pointer(here.tail)) = data[i]
		if here.tail == here.rightLimit {
			here.tail = here.leftLimit
		} else {
			here.tail += 1
		}
	}
	here.size += len(data)
	return nil
}

func (here *ShareMemory) ReadShareMemory(len int) []byte {
	if here.size == 0 {
		return make([]byte, 0)
	}
	ansLen := util.IntMin(len, here.size)
	ans := make([]byte, len)
	for i := 0; i < ansLen; i++ {
		ans[i] = *(*byte)(unsafe.Pointer(here.head))
		here.head += 1
	}
	here.size -= ansLen
	return ans
}

func (here *ShareMemory) Reset() {
	here.head = here.leftLimit
	here.tail = here.head
	here.size = 0
}
