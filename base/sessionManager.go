package base

import (
	"log"
	"sync"
)

type SessionManager struct {
	sessions map[int64]*Session
	lock     sync.RWMutex
}

func (self *SessionManager) AddSession(s *Session) {
	self.lock.Lock()
	defer self.lock.Unlock()

	if v, ok := self.sessions[s.ID()]; ok {
		log.Printf("session repeat %d", v.ID())
	}
	self.sessions[s.ID()] = s
}

func (self *SessionManager) RemoveSession(s *Session) {
	self.lock.Lock()
	defer self.lock.Unlock()

	if _, ok := self.sessions[s.ID()]; ok {
		delete(self.sessions, s.ID())
		log.Printf("session remove %d", s.ID())
	}
}

func (self *SessionManager) RemoveSessionByID(id int64) {
	self.lock.Lock()
	defer self.lock.Unlock()

	if _, ok := self.sessions[id]; ok {
		delete(self.sessions, id)
		log.Printf("session remove %d", id)
	}
}

func (self *SessionManager) GetSession(id int64) *Session {
	self.lock.RLock()
	defer self.lock.RUnlock()

	if v, ok := self.sessions[id]; ok {
		return v
	} else {
		return nil
	}
}

func NewSessionManager(n int) *SessionManager {
	m := &SessionManager{
		sessions: make(map[int64]*Session, n),
	}
	return m
}
