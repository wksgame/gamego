package base

import (
	"log"
	"net"
	"time"
)

type Server struct {
	netSrv     *NetServer
	Proc       *Processor
	exit       chan bool
	Manager    *SessionManager
	verifyFunc func(pkt *Packet) bool //检查连接是否合法
}

// 实现ConnectHandler接口
func (self *Server) SetNetServer(ns *NetServer) {
	self.netSrv = ns
}

// 实现ConnectHandler接口
func (self *Server) NewConnect(conn net.Conn) {
	go self.Verify(conn)
}

// 验证连接合法性
// 第一个包作为验证包
// 验证失败直接断开连接
func (self *Server) Verify(conn net.Conn) {
	stream := NewPacketStream(conn)
	for {
		select {
		case pkt, ok := <-stream.ReadChan():
			if !ok {
				log.Printf("Session verify failed, can't get packet")
				return
			}
			if self.verifyFunc(pkt) {
				log.Printf("Session verify ok")
				session := newSession(stream, self)
				session.SetID(int64(pkt.MsgID))
				self.Manager.AddSession(session)
				go session.Run()
			} else {
				stream.Close()
				log.Printf("Session verify failed, invalid token")
			}
			return
		case <-time.After(time.Second * 10):
			stream.Close()
			log.Printf("Session verify failed, timeout")
			return
		}
	}
}

func (self *Server) Stop(arg interface{}) {
	close(self.exit)
}

// 账号检查
func CheckAccount(pkt *Packet) bool {
	token := string(pkt.Data)
	return token == "hehe"
}

func NewServer(port int, pro *Processor) error {
	server := &Server{
		Proc:       pro,
		exit:       make(chan bool),
		Manager:    NewSessionManager(1000),
		verifyFunc: CheckAccount,
	}

	ExitApplication(server.Stop, nil)

	ns := &NetServer{}
	err := ns.ListenAndServe(port, server)
	return err
}
