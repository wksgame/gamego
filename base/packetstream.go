package base

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"log"
	"net"
)

type Sender interface {
}

type Packet struct {
	Sender Sender // 包的来源
	MsgID  int32  // 消息ID
	Data   []byte
}

const (
	PackageHeaderSize = 8 // DataLen(int32) + MsgID(int32)
	MaxPacketSize     = 1024 * 8
)

var (
	packageDataSizeInvalid = errors.New("ReadPacket: package invalid size")
	packageTooBig          = errors.New("ReadPacket: package too big")
)

// 封包流
type PacketStream interface {
	Close() error
	Raw() net.Conn

	// 调用此接口从流中读取数据
	// 例如：
	//	for {
	//		select {
	//		case pkt := <-stream.ReadChan():
	//			// todo
	//		}
	//	}
	ReadChan() <-chan *Packet

	// 调用此接口向流中写入数据
	WriteChan() chan<- *Packet
}

type packetStream struct {
	conn    net.Conn
	msghead []byte       // read message head
	r       chan *Packet // 内部从socket读取数据存入r,外部调用接口从r读取数据
	w       chan *Packet // 内部从w读取数据写入socket,外部调用接口写入数据到w
}

// 实现PacketStream接口
func (self *packetStream) Close() error {
	close(self.w)
	return self.conn.Close()
}

// 实现PacketStream接口
func (self *packetStream) Raw() net.Conn {
	return self.conn
}

// 实现PacketStream接口
func (self *packetStream) ReadChan() <-chan *Packet {
	return self.r
}

// 实现PacketStream接口
func (self *packetStream) WriteChan() chan<- *Packet {
	return self.w
}

func (self *packetStream) Go() {
	go self.readGo()
	go self.writeGo()
}

func (self *packetStream) readGo() {
	for {
		p, err := self.read()
		if err != nil {
			close(self.r)
			//self.conn.Close()

			log.Printf("packetStream recv error:%s", err)
			return
		}
		self.r <- p
	}
}

func (self *packetStream) writeGo() {
	ok := true
	for p := range self.w {
		if !ok {
			continue
		}
		err := self.write(p)
		if err != nil {
			log.Printf("writeGo exit, err:%s", err)
			ok = false
		}
	}
	log.Printf("wirteGo exit")
}

// 从socket读取数据
func (self *packetStream) read() (p *Packet, err error) {

	if _, err = io.ReadFull(self.conn, self.msghead); err != nil {
		return nil, err
	}

	p = &Packet{}

	headbuf := bytes.NewReader(self.msghead)

	// 读取消息长度
	var msgLen int32
	if err = binary.Read(headbuf, binary.LittleEndian, &msgLen); err != nil {
		return nil, err
	}

	// 读取消息ID
	if err = binary.Read(headbuf, binary.LittleEndian, &p.MsgID); err != nil {
		return nil, err
	}

	// 封包太大
	if msgLen > MaxPacketSize {
		return nil, packageTooBig
	}

	if msgLen < 0 {
		return nil, packageDataSizeInvalid
	}

	// 读取数据
	p.Data = make([]byte, msgLen)
	if _, err = io.ReadFull(self.conn, p.Data); err != nil {
		return nil, err
	}

	return
}

// 发送数据到socket
func (self *packetStream) write(pkt *Packet) (err error) {

	outbuff := bytes.NewBuffer([]byte{})

	// 写入消息长度
	if err = binary.Write(outbuff, binary.LittleEndian, int32(len(pkt.Data))); err != nil {
		return
	}

	// 写入消息ID
	if err = binary.Write(outbuff, binary.LittleEndian, pkt.MsgID); err != nil {
		return
	}

	// 写入数据
	if err = binary.Write(outbuff, binary.LittleEndian, pkt.Data); err != nil {
		return
	}

	// 发送数据
	if _, err = self.conn.Write(outbuff.Bytes()); err != nil {
		return
	}

	return
}

func NewPacketStream(conn net.Conn) PacketStream {
	stream := &packetStream{
		conn:    conn,
		msghead: make([]byte, PackageHeaderSize),
		r:       make(chan *Packet, 100),
		w:       make(chan *Packet, 100),
	}
	go stream.Go()
	return stream
}
