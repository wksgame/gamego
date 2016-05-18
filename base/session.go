package base

import (
	"log"
	"net"
)

type Session struct {
	stream PacketStream
	srv    *Server
	exit   chan bool
	id     int64
}

func (self *Session) ID() int64 {
	return self.id
}

func (self *Session) Close() {
	close(self.exit)
}

func (self *Session) Send(msgid int32, msg []byte) {
	pkt := &Packet{
		MsgID: msgid,
		Data:  msg,
	}
	self.stream.WriteChan() <- pkt
}

func (self *Session) Run() {
	R := self.stream.ReadChan()
	for {
		select {
		case pkt, ok := <-R:
			if ok {
				log.Printf("Session recv message")
				pkt.Sess = self
				self.srv.Proc.PushMessage(pkt)
			} else {
				return
			}
		case <-self.exit:
			return
		}
	}
}

func newSession(c net.Conn, s *Server) *Session {
	ses := &Session{
		stream: NewPacketStream(c),
		srv:    s,
		exit:   make(chan bool),
	}
	return ses
}
