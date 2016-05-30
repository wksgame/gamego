package base

import (
	"log"
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
	//R := self.stream.ReadChan()
	for {
		select {
		case pkt, ok := <-self.stream.ReadChan(): //<-R:
			if ok {
				log.Printf("Session recv message")
				pkt.Sender = self
				self.srv.Proc.PushMessage(pkt)
			} else {
				return
			}
		case <-self.exit:
			return
		}
	}
}

func newSession(stream PacketStream, s *Server) *Session {
	ses := &Session{
		stream: stream,
		srv:    s,
		exit:   make(chan bool),
	}
	return ses
}
