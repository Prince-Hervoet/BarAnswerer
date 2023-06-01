package main

import "ShareMemTCP/core"

func test1() {
	core.NewPoineer(20000)
}

func test2() {
	x := core.NewPoineer(20001)
	sid, _ := x.Link(20000, 4096)
	x.Send([]byte("asdfasdf"), sid)
}

func main() {
	// test1()
	// for {

	// }
	test2()
	for {

	}
}
