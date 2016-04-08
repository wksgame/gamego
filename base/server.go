package base

import (
	"log"
	"net"
	"strconv"
)

type Server struct {
	listener *net.TCPListener
	Proc     *Processor
	exit     chan bool
}

func (self *Server) listenAndServe(port int) error {
	serverAddr, err := net.ResolveTCPAddr("tcp4", ":"+strconv.Itoa(port))
	if err != nil {
		return err
	}

	l, err := net.ListenTCP("tcp", serverAddr)
	if err != nil {
		return err
	}

	self.listener = l

	return nil
}

func (self *Server) Start() {
	for {
		c, err := self.listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		s := newSession(c, self)

		go s.Run()
	}
}

func (self *Server) Stop(arg interface{}) {
	close(self.exit)
}

func NewServer(port int) (*Server, error) {
	server := &Server{
		Proc: newProcessor(),
		exit: make(chan bool),
	}

	err := server.listenAndServe(port)

	if err != nil {
		return nil, err
	} else {
		return server, nil
	}
}
