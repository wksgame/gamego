package base

import (
	"log"
	"net"
)

type Server struct {
	netSrv   *netServer
	Proc     *Processor
	exit     chan bool
	sessions map[int64]*Session
}

func (self *Server) Accept(conn net.Conn) {
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

func (self *Server) Start() {
	self.netSrv.Start(self.Accept)
}

func (self *Server) Stop(arg interface{}) {
	close(self.exit)
}

func NewServer(port int, pro *Processor) (*Server, error) {
	net, err := NewNetServer(port)

	if err != nil {
		return nil, err
	}

	server := &Server{
		netSrv:   net,
		Proc:     pro,
		exit:     make(chan bool),
		sessions: make(map[int64]*Session),
	}

	ExitApplication(server.Stop, nil)

	return server, nil
}
