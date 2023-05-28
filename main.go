package main

import (
	"ShareMemTCP/core"
	memory "ShareMemTCP/memory"
	"fmt"
	"time"
)

func main() {
	s := &core.Sharer{}
	s.Open("asdf")
}

func test1() {
	sm, _, err := memory.OpenShareMemory(4096)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	test := make([]byte, 1)
	for i := 0; i < len(test); i++ {
		test[i] = byte(i)
	}
	sm.WriteShareMemory(test)

	for true {
		fmt.Println("wait...")
		time.Sleep(2 * time.Second)
	}
}

func test2() {
	memory.GetShareMemory(1, 4096)
}
