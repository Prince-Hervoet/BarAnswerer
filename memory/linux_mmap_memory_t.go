package memory

import (
	"fmt"
)

func Test1() {
	msm := NewMmapShareMemory()
	msm.OpenFile("/test.txt", 2048)
	msm.Mmap()
	// mh := protocol.NewMemoryHeader(1, 123, 333)
	// msm.WriteHeader(mh)
	mh := msm.ReadHeader()
	fmt.Println(mh)
	msm.CancelMmap()
}

func Test2() {
	msm := NewMmapShareMemory()
	msm.OpenFile("/test.txt", 2048)
	msm.Mmap()
	test := make([]byte, 10)
	for i := 0; i < len(test); i++ {
		test[i] = byte(i)
	}
	msm.WriteBytes(test)
	mh := msm.ReadBytes(10)
	mh2 := msm.ReadHeader()
	fmt.Println(mh)
	fmt.Println(mh2)
}
