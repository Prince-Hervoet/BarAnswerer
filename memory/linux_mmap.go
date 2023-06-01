package memory

import (
	"ShareMemTCP/util"
	"errors"
	"fmt"
	"os"
	"syscall"
)

type ShareMemory struct {
	filePath string
	filePtr  *os.File
	isOpened bool
	mapping  []byte
	header   ShareMemoryHeader
}

func OpenShareMemory() *ShareMemory {
	return &ShareMemory{
		filePath: "",
		filePtr:  nil,
		isOpened: false,
	}
}

func (here *ShareMemory) Size() int32 {
	here.readHeader()
	return here.header.size
}

func (here *ShareMemory) OpenFile(fileName string, cap int32) (string, error) {
	if here.isOpened {
		return "", errors.New("a mapping has been established")
	} else if cap <= util.SHARE_MEMORY_HEADER_SIZE {
		return "", errors.New("cap is too small")
	}
	finalPath := util.DEFAULT_FILE_DIR + fileName
	if !PathExists(util.DEFAULT_FILE_DIR) {
		os.Mkdir(util.DEFAULT_FILE_DIR, os.ModePerm)
	}
	filePtr, err := os.OpenFile(finalPath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Println("open file error")
		filePtr.Close()
		return "", err
	}
	here.filePtr = filePtr
	here.header.cap = cap
	bs, err := here.mmap()
	if err != nil {
		fmt.Println("mmap file error")
		filePtr.Close()
		return "", err
	}
	here.grow()
	here.filePath = finalPath
	here.isOpened = true
	here.mapping = bs
	here.header.size = util.SHARE_MEMORY_HEADER_SIZE
	here.header.status = util.NO_WORKING
	here.header.head = util.SHARE_MEMORY_HEADER_SIZE + 1
	here.header.tail = util.SHARE_MEMORY_HEADER_SIZE + 1
	here.writeHeader()
	return finalPath, nil
}

// 连接已经有的文件
func (here *ShareMemory) LinkFile(filePath string, cap int32) error {
	if here.isOpened {
		return errors.New("a mapping has been established")
	} else if cap <= util.SHARE_MEMORY_HEADER_SIZE {
		return errors.New("cap is too small")
	}
	filePtr, err := os.OpenFile(filePath, os.O_RDWR, 0666)
	if err != nil {
		fmt.Println(err.Error())
		filePtr.Close()
		return err
	}
	here.filePtr = filePtr
	here.header.cap = cap
	bs, err := here.mmap()
	if err != nil {
		fmt.Println("mmap file error")
		filePtr.Close()
		return err
	}
	here.filePath = filePath
	here.isOpened = true
	here.mapping = bs
	here.readHeader()
	return nil
}

func (here *ShareMemory) Close() error {
	if !here.isOpened {
		return nil
	}
	err := here.munmap()
	if err != nil {
		return err
	}
	err = here.filePtr.Close()
	if err != nil {
		return err
	}
	os.Remove(here.filePath)
	here.filePtr = nil
	here.filePath = ""
	here.mapping = nil
	here.isOpened = false
	here.header.cap = 0
	here.header.size = 0
	here.header.status = 0
	here.header.head = 0
	here.header.tail = 0
	return nil
}

func (here *ShareMemory) Write(data []byte) error {
	if !here.isOpened {
		return errors.New("please open a file")
	} else if len(data)+int(here.header.size) > int(here.header.cap) {
		return errors.New("no enough sapce")
	}

	tail := here.header.tail
	for i := 0; i < len(data); i++ {
		here.mapping[tail] = data[i]
		if tail == here.header.cap-1 {
			tail = util.SHARE_MEMORY_HEADER_SIZE + 1
		} else {
			tail += 1
		}
	}
	here.header.tail = tail
	here.header.size += int32(len(data))
	here.writeHeader()
	return nil
}

func (here *ShareMemory) Read(bs []byte) (int32, error) {
	if !here.isOpened {
		return 0, errors.New("please open a file")
	}
	here.readHeader()
	if here.header.size == 0 {
		return 0, nil
	}
	ansLen := util.IntMin(len(bs), int(here.header.size))
	head := here.header.head
	for i := 0; i < ansLen; i++ {
		bs[i] = here.mapping[head]
		if head == here.header.cap-1 {
			head = util.SHARE_MEMORY_HEADER_SIZE + 1
		} else {
			head += 1
		}
	}
	here.header.head = head
	here.header.size -= int32(ansLen)
	here.writeHeader()
	return int32(ansLen), nil
}

func (here *ShareMemory) ChangeStatus(status int8) {
	if !here.isOpened {
		return
	}
	here.mapping[0] = byte(status)
	here.writeHeader()
}

func (here *ShareMemory) mmap() ([]byte, error) {
	bs, err := syscall.Mmap(int(here.filePtr.Fd()), 0, int(here.header.cap), syscall.PROT_WRITE|syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return bs, nil
}

func (here *ShareMemory) munmap() error {
	err := syscall.Munmap(here.mapping)
	if err != nil {
		fmt.Println("close error")
		return err
	}
	return nil
}

func (here *ShareMemory) grow() error {
	err := here.filePtr.Truncate(int64(here.header.cap))
	if err != nil {
		return err
	}
	return nil
}

func (here *ShareMemory) readHeader() {
	if !here.isOpened {
		return
	}
	temp := here.mapping[0 : util.SHARE_MEMORY_HEADER_SIZE+1]
	here.header.FromBytes(temp)
}

func (here *ShareMemory) writeHeader() {
	if !here.isOpened {
		return
	}
	bs := here.header.ToBytes()
	for i := 0; i < len(bs); i++ {
		here.mapping[i] = bs[i]
	}
}

func (here *ShareMemory) checkStatus() bool {
	if !here.isOpened {
		return false
	}
	if here.mapping[0] == byte(0) {
		return false
	} else if here.mapping[0] == byte(1) {
		return true
	}
	return false
}

func PathExists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		return os.IsExist(err)
	}
	return true
}
