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

// 用于记录文件打开次数
var fileNameSeq int = 0

const (
	// 用于存储文件的路径
	DEFAULT_TMP_PATH = "/tmp/share"
)

// 映射内存数据结构
type MmapShareMemory struct {
	bufferHeader *protocol.MemoryHeader
	filePath     string
	file         *os.File
	memoryPtr    []byte
	startPtr     uintptr
	cap          int32
	size         int32
	head         int32
	tail         int32
}

func NewMmapShareMemory() *MmapShareMemory {
	return &MmapShareMemory{
		filePath:     "",
		file:         nil,
		size:         util.MEM_HEADER_SIZE,
		head:         util.MEM_HEADER_SIZE + 1,
		tail:         util.MEM_HEADER_SIZE + 1,
		cap:          0,
		bufferHeader: protocol.NewMemoryHeader(0, 0, 0),
	}
}

// 打开系统文件，准备作为mmap使用
func (here *MmapShareMemory) OpenFile(filePath string, cap int32) error {
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
	here.size = util.MEM_HEADER_SIZE
	here.file = file
	here.filePath = finalPath
	here.grow(cap)
	fileNameSeq += 1
	return nil
}

// 将file指针指向的文件进行一个映射
func (here *MmapShareMemory) Mmap() error {
	bs, err := syscall.Mmap(int(here.file.Fd()), 0, int(here.cap), syscall.PROT_WRITE|syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		fmt.Println("mmap error")
		return err
	}
	here.memoryPtr = bs
	here.startPtr = uintptr(unsafe.Pointer(&bs[0]))
	return nil
}

// 设置内部文件的大小
func (here *MmapShareMemory) grow(size int32) {
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
	if here.memoryPtr == nil {
		return nil
	}
	err := syscall.Munmap(here.memoryPtr)
	if err != nil {
		return err
	}
	here.cap = 0
	here.size = 0
	here.startPtr = 0
	here.memoryPtr = nil
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
func (here *MmapShareMemory) ReadHeader() *protocol.MemoryHeader {
	temp := here.memoryPtr[0 : util.MEM_HEADER_SIZE+1]
	here.bufferHeader.FromByteArray(temp)
	here.head = here.bufferHeader.Head
	here.tail = here.bufferHeader.Tail
	return here.bufferHeader
}

func (here *MmapShareMemory) WriteHeader(header *protocol.MemoryHeader) {
	ans := header.ToByteArray()
	for i := 0; i < len(ans); i++ {
		here.memoryPtr[i] = ans[i]
	}
}

func (here *MmapShareMemory) WriteBytes(data []byte) error {
	if int32(len(data))+here.size > here.cap {
		return errors.New("no enough space")
	}
	for i := 0; i < len(data); i++ {
		here.memoryPtr[here.tail] = data[i]
		if here.tail == here.cap-1 {
			here.tail = util.MEM_HEADER_SIZE + 1
		} else {
			here.tail += 1
		}
	}
	here.size += int32(len(data))
	here.bufferHeader.Flag = util.MEM_FLAG_WORKING
	here.bufferHeader.Head = here.head
	here.bufferHeader.Tail = here.tail
	here.WriteHeader(here.bufferHeader)
	return nil
}

func (here *MmapShareMemory) ReadBytes(len int) []byte {
	if here.file == nil {
		return nil
	} else if here.size == 0 {
		return make([]byte, 0)
	}
	ansLen := util.IntMin(len, int(here.size))
	ans := make([]byte, ansLen)
	for i := 0; i < ansLen; i++ {
		ans[i] = here.memoryPtr[here.head]
		if here.head == here.cap-1 {
			here.head = util.MEM_HEADER_SIZE + 1
		} else {
			here.head += 1
		}
	}
	here.size -= int32(ansLen)
	here.bufferHeader.Flag = util.MEM_FLAG_COMMON
	here.bufferHeader.Head = here.head
	here.bufferHeader.Tail = here.tail
	here.WriteHeader(here.bufferHeader)
	return ans
}
