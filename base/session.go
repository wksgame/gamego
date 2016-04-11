package base

import (
	"log"
	"net"

	//	"github.com/golang/protobuf/proto"
)

type Session struct {
	writeChan chan *Packet
	stream    PacketStream
	srv       *Server
	exit      chan bool

	OnClose func() // 关闭函数回调

	id int64
}

func (self *Session) ID() int64 {
	return self.id
}

func (self *Session) Run() {

	// 接收线程
	go self.recvThread()

	// 发送线程
	go self.sendThread()
}

func (self *Session) Close() {
	close(self.exit)
}

func (self *Session) Send(msgid int32, msg []byte) {

	//	data, err := proto.Marshal(msg)

	//	if err != nil {
	//		log.Println("send err,", err)
	//		return
	//	}

	pkt := &Packet{
		MsgID: msgid,
		Data:  msg,
	}

	self.RawSend(pkt)
}

func (self *Session) RawSend(pkt *Packet) {

	if pkt == nil {
		return
	}

	self.writeChan <- pkt
}

// 发送线程
func (self *Session) sendThread() {

	for {
		select {
		// 封包
		case pkt := <-self.writeChan:
			if err := self.stream.Write(pkt); err != nil {
				goto exitsendloop
			}
		case <-self.srv.exit:
			goto exitsendloop
		case <-self.exit:
			goto exitsendloop
		}

	}

exitsendloop:
	self.stream.Close()
}

// 接收线程
func (self *Session) recvThread() {
	var err error
	var pkt *Packet

	for {
		pkt, err = self.stream.Read()

		if err != nil {
			break
		}

		pkt.Sess = self
		self.srv.Proc.PushMessage(pkt)
	}

	self.Close()
}

func newSession(c net.Conn, s *Server) *Session {

	ses := &Session{
		writeChan: make(chan *Packet, 100),
		stream:    NewPacketStream(c),
		srv:       s,
		exit:      make(chan bool),
	}

	return ses
}
