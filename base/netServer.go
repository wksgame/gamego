package base

import (
	"log"
	"net"
	"strconv"
	"time"
)

type ConnectHandler interface {
	// 设置当前对应的NetServer
	SetNetServer(*NetServer)
	// 有新连接时的处理
	NewConnect(conn net.Conn)
}

type NetServer struct {
	listener *net.TCPListener
}

func (self *NetServer) ListenAndServe(port int, handler ConnectHandler) error {
	handler.SetNetServer(self)

	serverAddr, err := net.ResolveTCPAddr("tcp4", ":"+strconv.Itoa(port))
	if err != nil {
		return err
	}

	l, err := net.ListenTCP("tcp", serverAddr)
	if err != nil {
		return err
	}

	self.listener = l
	var tempDelay time.Duration // how long to sleep on accept failure

	for {
		conn, err := self.listener.Accept()

		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Second
				} else {
					tempDelay *= 2
				}
				if max := 60 * time.Second; tempDelay > max {
					tempDelay = max
				}
				log.Printf("NetServer Accept error: %s; retrying in %s", err, tempDelay)
				time.Sleep(tempDelay)
				continue
			}
			return err
		}
		tempDelay = 0

		handler.NewConnect(conn)
	}

	return nil
}
