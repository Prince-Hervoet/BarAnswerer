package main

import "ShareMemTCP/core"

func test1() {
	x := core.Pioneer{}
	x.ConnectInit(20000)
}

func test2() {
	x := core.Pioneer{}
	x.ConnectInit(20001)
	seesionId, _ := x.OpenConnection(20000, 4096)
	x.Close(seesionId)
}

func main() {
	// test1()
	// for {

	// }
	test2()
}
