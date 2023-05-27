package main

import (
	"ShareMemTCP/core"
	"fmt"
)

func main() {
	sm := core.OpenShareMemory(21, 4096)
	test := make([]byte, 10)
	for i := 0; i < 10; i++ {
		test[i] = byte(i)
	}
	err := sm.WriteShareMemory(test)
	ans := sm.ReadShareMemory(10)
	fmt.Print("ans: ")
	for i := 0; i < len(ans); i++ {
		fmt.Print(ans[i])
		fmt.Print(" ")
	}
	fmt.Println()
	if err != nil {
		fmt.Println(err.Error())
	}
}
