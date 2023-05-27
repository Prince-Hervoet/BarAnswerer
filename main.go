package main

import (
	"ShareMemTCP/core"
	"fmt"
	"time"
)

func main() {

}

func test1() {
	sm, _, err := core.OpenShareMemory(4096)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	test := make([]byte, 10)
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
	core.GetShareMemory(1, 4096)
}
