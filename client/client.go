package main

import (
	"log"
	"net"
	"strconv"
	"time"

	. "../base"
)

var ipport string = "192.168.1.34:4444"

func ConnectServer(i int) {
	addr, err := net.ResolveTCPAddr("tcp4", ipport)
	if err != nil {
		log.Println("create addr err:", err)
		return
	}

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		log.Println("conn err:", err)
		return
	}

	stream := NewPacketStream(conn)

	p := &Packet{
		MsgID: 1,
		Data:  []byte("hehehehehe " + strconv.Itoa(i)),
	}

	for {
		stream.Write(p)
		pkt, _ := stream.Read()
		log.Println(string(pkt.Data))
		time.Sleep(time.Millisecond * 5)
	}
	log.Println("exit:", i)
}

func main() {
	//time.Sleep(time.Second * 3)

	for i := 0; i < 1; i++ {
		go ConnectServer(i)
		//		if i%50 == 0 {
		//			time.Sleep(time.Second)
		//		}
		//time.Sleep(time.Millisecond * 50)
		log.Println(i)
	}

	time.Sleep(time.Second * 3600)
}
