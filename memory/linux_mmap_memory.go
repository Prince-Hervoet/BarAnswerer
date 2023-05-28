package memory

import (
	"ShareMemTCP/protocol"
	"ShareMemTCP/util"
	"errors"
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

var fileNameSeq int = 1

const (
	DEFAULT_TMP_PATH = "/tmp"
)

type MmapShareMemory struct {
	filePath   string
	file       *os.File
	defaultPtr []byte
	startPtr   uintptr
	cap        int
	size       int
	sendQueue  *ActionSection
	getQueue   *ActionSection
}

func NewMmapShareMemory() *MmapShareMemory {
	return &MmapShareMemory{
		filePath: "",
		file:     nil,
		size:     0,
		cap:      0,
	}
}

// 打开系统文件，准备作为mmap使用
func (here *MmapShareMemory) OpenFile(filePath string, cap int) error {
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
		fmt.Println(err.Error())
		return err
	}
	here.cap = cap
	here.size = 0
	here.file = file
	here.Grow(cap)
	return nil
}

// 将file指针指向的文件进行一个映射
func (here *MmapShareMemory) RunMmap() error {
	bs, err := syscall.Mmap(int(here.file.Fd()), 0, here.cap, syscall.PROT_WRITE|syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		fmt.Println("mmap error")
		return err
	}
	here.defaultPtr = bs
	here.startPtr = uintptr(unsafe.Pointer(&bs[0]))
	return nil
}

// 设置内部文件的大小
func (here *MmapShareMemory) Grow(size int) {
	if info, _ := here.file.Stat(); info.Size() >= int64(size) {
		return
	}
	err := here.file.Truncate(int64(size))
	if err != nil {
		fmt.Println(err.Error())
	}
}

// 关闭映射（同时会关闭和删除文件）
func (here *MmapShareMemory) CancelMmap() error {
	if here.defaultPtr == nil {
		return nil
	}
	err := syscall.Munmap(here.defaultPtr)
	if err != nil {
		return err
	}
	here.cap = 0
	here.size = 0
	here.startPtr = 0
	here.defaultPtr = nil
	err = here.file.Close()
	here.file = nil
	if err != nil {
		return err
	}
	err = os.Remove(here.filePath)
	here.filePath = ""
	if err != nil {
		return err
	}
	return nil
}

// 获取共享内存协议头部
func (here *MmapShareMemory) ReadHeader() *protocol.MemoryHeaderProtocol {
	header := &protocol.MemoryHeaderProtocol{}
	temp := here.defaultPtr[0:10]
	header.FromByteArray(temp)
	return header
}

func (here *MmapShareMemory) WriteHeader(header *protocol.MemoryHeaderProtocol) {
	ans := header.ToByteArray()
	for i := 0; i < len(ans); i++ {
		here.defaultPtr[i] = ans[i]
	}
}

func (here *MmapShareMemory) WriteBytes(data []byte) error {
	if len(data)+here.size > here.cap {
		return errors.New("no enough space")
	}
	for i := 0; i < len(data); i++ {
		here.defaultPtr[i] = data[i]
		here.startPtr += 1
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
	fmt.Println(here.size)
	ansLen := util.IntMin(len, here.size)
	ans := make([]byte, ansLen)
	here.startPtr = uintptr(unsafe.Pointer(&here.defaultPtr[0]))
	for i := 0; i < ansLen; i++ {
		ans[i] = *(*byte)(unsafe.Pointer(here.startPtr))
		here.startPtr += 1
	}
	return ans
}

// 重置共享内存区域，会清空所有数据
func (here *MmapShareMemory) Reset() {
	if here.file != nil {
		return
	}
	here.startPtr = uintptr(unsafe.Pointer(&here.defaultPtr[0]))
	here.size = 0
}
