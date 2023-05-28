package core

import (
	"ShareMemTCP/connection"
	"ShareMemTCP/memory"
	"fmt"
)

type Sharer struct {
	mineMemories  []int
	otherMemories []int
	pioneer       *connection.Pioneer
}

func (here *Sharer) Open(address string) {
	msm := memory.NewMmapShareMemory()
	msm.OpenFile("/share/test.txt", 4096)
	msm.Grow(1)
	err := msm.RunMmap()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	test := make([]byte, 10)
	for i := 0; i < 10; i++ {
		test[i] = byte(i)
	}
	msm.Grow(len(test))
	msm.WriteBytes(test)

	ans := msm.ReadBytes(10)
	for i := 0; i < len(ans); i++ {
		fmt.Print(ans[i])
		fmt.Print(" ")
	}
	fmt.Println()
}

func (here *Sharer) Link(filePath string) {

}

func (here *Sharer) Send() {

}

func (here *Sharer) Read() {

}
