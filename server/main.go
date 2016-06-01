package main

import (
	"runtime"

	. "../base"
)

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	ini := &IniConfig{}
	err := ini.Parse("server.ini")
	if err != nil {
		return
	}
	port, err := ini.GetInt("server", "port")
	if err != nil {
		return
	}

	pro := NewProcessor()
	pro.RegisterMessage(0, true, OnLogin)
	pro.RegisterMessage(1, false, OnLogin)
	pro.RegisterMessage(2, true, OnLogout)

	NewServer(port, pro)
}
