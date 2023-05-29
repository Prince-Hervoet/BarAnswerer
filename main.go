package main

import "ShareMemTCP/connection"

func test1(){
	x := connection.Pioneer{}
	x.ConnectInit(2,":20001")
}

func test2(){
	x := connection.Pioneer{}
	x.OpenConnection(":20001", 4096)
}

func main() {
	// test1()
	// for {
		
	// }
	test2()
}
