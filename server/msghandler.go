package main

import (
	"log"

	. "../base"
)

func OnLogin(msg *Packet) {
	log.Println("login data:", string(msg.Data))
	//log.Println("login msgid:", msg.MsgID)

	msg.Sess.Send(msg.MsgID, msg.Data)
}
