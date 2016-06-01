package base

import (
	"log"
	"sync"
)

type Session struct {
	stream PacketStream
	srv    *Server
	exit   chan bool
	sendOk bool
	lock   sync.RWMutex
	id     int64
}

func (self *Session) ID() int64 {
	return self.id
}

func (self *Session) SetID(id int64) {
	self.id = id
}

func (self *Session) Close() {
	close(self.exit)
}

func (self *Session) Send(msgid int32, msg []byte) {
	self.lock.RLock()
	defer self.lock.RUnlock()

	if !self.sendOk {
		return
	}
	pkt := &Packet{
		MsgID: msgid,
		Data:  msg,
	}
	self.stream.WriteChan() <- pkt
}

func (self *Session) Run() {
	for {
		select {
		case pkt, ok := <-self.stream.ReadChan():
			if !ok {
				log.Printf("Session recv error")
				goto quit
			}
			log.Printf("Session recv message")
			pkt.Sender = self
			self.srv.Proc.PushMessage(pkt)
		case <-self.exit:
			log.Printf("Session recv quit")
			goto quit
		}
	}
quit:
	self.lock.Lock()
	self.sendOk = false
	self.stream.Close()
	self.lock.Unlock()

	self.srv.Manager.RemoveSession(self)
}

func newSession(stream PacketStream, s *Server) *Session {
	ses := &Session{
		stream: stream,
		srv:    s,
		exit:   make(chan bool),
		sendOk: true,
	}
	return ses
}
