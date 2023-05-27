package core

import (
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

func OpenShareMemory(cap int) *ShareMemory {
	shmid, _, err := syscall.Syscall(syscall.SYS_SHMGET, 1234, 12, IPC_CREATE|0666)
	if err != 0 {
		fmt.Println("open memory error")
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
		*(*byte)(unsafe.Pointer(here.head)) = data[i]
		here.head += 1
	}
	here.size += len(data)
	fmt.Println(here.size)
	return nil
}

func (here *ShareMemory) Reset() {

}

func ReadShareMemory() {

}
