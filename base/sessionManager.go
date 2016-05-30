package base

import (
	"log"
)

type SessionManager struct {
	sessions map[int64]*Session
}

func (self *SessionManager) AddSession(s *Session) {
	if v, ok := self.sessions[s.ID()]; ok {
		log.Println("session repeat", v.ID())
	}
	self.sessions[s.ID()] = s
}

func (self *SessionManager) RemoveSession(s *Session) {
	if _, ok := self.sessions[s.ID()]; ok {
		delete(self.sessions, s.ID())
		log.Println("session online")
	}
}

func (self *SessionManager) RemoveSessionByID(id int64) {
	if _, ok := self.sessions[id]; ok {
		delete(self.sessions, id)
		log.Println("session online")
	}
}
