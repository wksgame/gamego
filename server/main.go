package main

import (
	"runtime"

	. "../base"
)

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	pro := NewProcessor()
	pro.RegisterMessage(0, true, OnLogin)
	pro.RegisterMessage(1, false, OnLogin)

	NewServer(4444, pro)
}
