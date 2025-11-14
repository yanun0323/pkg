package main

import "github.com/yanun0323/pkg/sys"

func main() {
	<-sys.Shutdown()
	println("AAA")

	<-sys.Shutdown()
	println("BBB")
}
