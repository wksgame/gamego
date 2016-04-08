package base

import (
	"log"
)

type MsgCallBack func(msg *Packet)

type MsgProc struct {
	runnow   bool
	callback MsgCallBack
}

type Processor struct {
	callbackMap map[int32]MsgProc
	msgqueue    chan *Packet
}

func (self *Processor) PushMessage(msg *Packet) {
	if cb, ok := self.callbackMap[msg.MsgID]; ok {
		if cb.runnow {
			cb.callback(msg)
		} else {
			self.msgqueue <- msg
		}
	} else {
		log.Printf("Processor.PushMessage err, msgid:%d", msg.MsgID)
	}
}

func (self *Processor) RegisterMessage(msgID int32, runnow bool, cb MsgCallBack) {
	if _, ok := self.callbackMap[msgID]; ok {
		log.Printf("Processor.RegisterMessage err, msgid:%d", msgID)
	}

	self.callbackMap[msgID] = MsgProc{runnow, cb}
}

func (self *Processor) ProcessMessage() {
	for {
		select {
		case msg := <-self.msgqueue:
			if cb, ok := self.callbackMap[msg.MsgID]; ok {
				cb.callback(msg)
				log.Println("process message ok")
			}
		}
	}
}

func newProcessor() *Processor {
	p := &Processor{
		callbackMap: make(map[int32]MsgProc),
		msgqueue:    make(chan *Packet, 1000000),
	}

	go p.ProcessMessage()
	return p
}
