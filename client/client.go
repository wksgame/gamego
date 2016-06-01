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

	log.Printf("connect ok id:%d", i)

	stream := NewPacketStream(conn)

	//time.Sleep(time.Second * 15)

	vp := &Packet{
		MsgID: int32(i),
		Data:  []byte("hehe"),
	}

	p := &Packet{
		MsgID: 1,
		Data:  []byte("message send by " + strconv.Itoa(i)),
	}

	q := &Packet{
		MsgID: 2,
		Data:  []byte("exit" + strconv.Itoa(i)),
	}

	rc := stream.ReadChan()
	wc := stream.WriteChan()

	wc <- vp
	loop := 0

	for {
		select {
		case pkt, ok := <-rc:
			if !ok {
				goto exit
			}
			log.Println(string(pkt.Data))
		default:
			wc <- p
			log.Printf("send message")
			time.Sleep(time.Millisecond * 500)
			loop++
			if loop > 5 {
				wc <- q
				goto exit
			}
		}
	}

exit:
	log.Println("exit:", i)
}

func main() {
	ini := &IniConfig{}
	if err := ini.Parse("client.ini"); err != nil {
		return
	}
	ip, err := ini.GetValue("client", "ip")
	if err != nil {
		return
	}
	port, err := ini.GetValue("client", "port")
	if err != nil {
		return
	}
	ipport = ip + ":" + port
	//time.Sleep(time.Second * 3)

	for i := 0; i < 10; i++ {
		go ConnectServer(i)
	}

	time.Sleep(time.Second * 3600)
}
