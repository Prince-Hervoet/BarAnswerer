package core

import (
	"ShareMemTCP/util"
	"errors"
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

const (
	DEFAULT_TMP_PATH = "/tmp"
)

type MmapShareMemory struct {
	file       *os.File
	defaultPtr []byte
	dataPtr    uintptr
	cap        int
	size       int
}

func NewMmapShareMemory() *MmapShareMemory {
	return &MmapShareMemory{
		file: nil,
	}
}

func (here *MmapShareMemory) Init(filePath string, cap int) error {
	if here.file != nil {
		err := here.file.Close()
		if err != nil {
			fmt.Println("close error")
			return err
		}
		here.file = nil
	}
	finalPath := DEFAULT_TMP_PATH + filePath
	file, err := os.OpenFile(finalPath, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("open file error")
		return err
	}
	here.cap = cap
	here.size = 0
	here.file = file
	return nil
}

func (here *MmapShareMemory) RunMmap() error {
	bs, err := syscall.Mmap(int(here.file.Fd()), 0, here.cap, syscall.PROT_WRITE|syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		fmt.Println("mmap error")
		return err
	}
	here.defaultPtr = bs
	here.dataPtr = uintptr(unsafe.Pointer(&bs[0]))
	return nil
}

func (here *MmapShareMemory) CancelMmap() error {
	if here.defaultPtr == nil {
		return nil
	}
	err := syscall.Munmap(here.defaultPtr)
	if err != nil {
		return err
	}
	err = here.file.Close()
	here.file = nil
	if err != nil {
		return err
	}
	here.cap = 0
	here.size = 0
	here.dataPtr = 0
	here.defaultPtr = nil
	return nil
}

func (here *MmapShareMemory) WriteBytes(data []byte) error {
	if len(data)+here.size > here.cap {
		return errors.New("no enough space")
	}
	for i := 0; i < len(data); i++ {
		*(*byte)(unsafe.Pointer(here.dataPtr)) = data[i]
		here.dataPtr += 1
	}
	here.size += len(data)
	return nil
}

func (here *MmapShareMemory) ReadBytes(len int) []byte {
	if here.file == nil {
		return nil
	} else if here.size == 0 {
		return make([]byte, 0)
	}
	ansLen := util.IntMin(len, here.size)
	ans := make([]byte, ansLen)
	for i := 0; i < ansLen; i++ {
		ans[i] = *(*byte)(unsafe.Pointer(here.dataPtr))
	}
	return nil
}

func (here *MmapShareMemory) Reset() {
	if here.file != nil {
		return
	}
	here.dataPtr = uintptr(unsafe.Pointer(&here.defaultPtr[0]))
	here.size = 0
}
