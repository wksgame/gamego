package main

import (
	"runtime"

	. "../base"
)

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	s, _ := NewServer(4444)

	s.Proc.RegisterMessage(0, true, OnLogin)
	s.Proc.RegisterMessage(1, false, OnLogin)

	ExitApp(s.Stop, nil)

	s.Start()
}
