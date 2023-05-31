package main

import "ShareMemTCP/core"

func test1() {
	x := core.Pioneer{}
	x.ConnectInit(20000)
}

func test2() {
	x := core.Pioneer{}
	x.ConnectInit(20000)
	x.Client.Link(20000, 4096)
}

func main() {
	// test1()
	// for {

	// }
	test2()
}
