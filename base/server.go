package base

import (
	"log"
	"net"
)

type Server struct {
	netSrv   *NetServer
	Proc     *Processor
	exit     chan bool
	sessions map[int64]*Session
}

// 实现ConnectHandler接口
func (self *Server) SetNetServer(ns *NetServer) {
	self.netSrv = ns
}

// 实现ConnectHandler接口
func (self *Server) NewConnect(conn net.Conn) {
	s := newSession(conn, self)
	go s.Run()
}

func (self *Server) AddSession(s *Session) {
	if v, ok := self.sessions[s.ID()]; ok {
		log.Println("session repeat", v.ID())
	}
	self.sessions[s.ID()] = s
}

func (self *Server) RemoveSession(s *Session) {
	if _, ok := self.sessions[s.ID()]; ok {
		delete(self.sessions, s.ID())
		log.Println("session online")
	}
}

func (self *Server) RemoveSessionByID(id int64) {
	if _, ok := self.sessions[id]; ok {
		delete(self.sessions, id)
		log.Println("session online")
	}
}

func (self *Server) Stop(arg interface{}) {
	close(self.exit)
}

func NewServer(port int, pro *Processor) error {
	server := &Server{
		Proc:     pro,
		exit:     make(chan bool),
		sessions: make(map[int64]*Session),
	}

	ExitApplication(server.Stop, nil)

	ns := &NetServer{}
	err := ns.ListenAndServe(port, server)
	return err
}
