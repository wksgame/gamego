package base

import (
	"log"
	"net"
	"strconv"
)

type Client struct {
	addr   *net.TCPAddr
	stream PacketStream
	Proc   *Processor
	exit   chan bool
}

func (self *Client) Create(ip string, port int) {
	ipport := ip + ":" + strconv.Itoa(port)
	addr, err := net.ResolveTCPAddr("tcp4", ipport)
	if err != nil {
		log.Println("create addr err:", err)
		return
	}
	self.addr = addr
}

func (self *Client) Connect() {
	conn, err := net.DialTCP("tcp", nil, self.addr)
	if err != nil {
		log.Println("conn err:", err)
		return
	}

	self.stream = NewPacketStream(conn)
}

func (self *Client) Send(msgid int32, msg []byte) {
	pkt := &Packet{
		MsgID: msgid,
		Data:  msg,
	}
	self.stream.WriteChan() <- pkt
}

func (self *Client) Run() {
	R := self.stream.ReadChan()
	for {
		select {
		case pkt, ok := <-R:
			if ok {
				log.Printf("Client recv message")
				pkt.Sender = self
				self.Proc.PushMessage(pkt)
			} else {
				return
			}
		case <-self.exit:
			return
		}
	}
}

func NewClient(ip string, port int, proc *Processor) *Client {
	client := &Client{
		Proc: proc,
		exit: make(chan bool),
	}
	return client
}
