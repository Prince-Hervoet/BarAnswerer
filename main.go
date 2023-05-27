package main

import (
	"ShareMemTCP/core"
)

func main() {
	sm := core.OpenShareMemory(4096)
	test := make([]byte, 5)
	sm.WriteShareMemory(test)
}
