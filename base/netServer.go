package base

import (
	"log"
	"net"
	"strconv"
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

	for {
		conn, err := self.listener.Accept()
		if err != nil {
			log.Printf("NetServer accept error:%s", err)
			continue
		}

		handler.NewConnect(conn)
	}

	return nil
}
