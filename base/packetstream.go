package base

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"net"
	"sync"
)

type Packet struct {
	Sess  *Session
	MsgID int32 // 消息ID
	Data  []byte
}

const (
	PackageHeaderSize = 8 // DataLen(int32) + MsgID(int32)
	MaxPacketSize     = 1024 * 8
)

// 封包流
type PacketStream interface {
	Read() (*Packet, error)
	Write(pkt *Packet) error
	Close() error
	Raw() net.Conn
}

type packetStream struct {
	conn         net.Conn
	sendtagGuard sync.RWMutex
	msghead      []byte // read message head
}

var (
	packageTagNotMatch     = errors.New("ReadPacket: package tag not match")
	packageDataSizeInvalid = errors.New("ReadPacket: package crack, invalid size")
	packageTooBig          = errors.New("ReadPacket: package too big")
)

// 从socket读取1个封包,并返回
func (self *packetStream) Read() (p *Packet, err error) {

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

// 将一个封包发送到socket
func (self *packetStream) Write(pkt *Packet) (err error) {

	outbuff := bytes.NewBuffer([]byte{})

	// 防止将Send放在go内造成的多线程冲突问题
	self.sendtagGuard.Lock()
	defer self.sendtagGuard.Unlock()

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

	// 发包头
	if _, err = self.conn.Write(outbuff.Bytes()); err != nil {
		return err
	}

	return
}

// 关闭
func (self *packetStream) Close() error {
	return self.conn.Close()
}

// 裸socket操作
func (self *packetStream) Raw() net.Conn {
	return self.conn
}

// 封包流 relay模式: 在封包头有clientid信息
func NewPacketStream(conn net.Conn) PacketStream {
	return &packetStream{
		conn:    conn,
		msghead: make([]byte, PackageHeaderSize),
	}
}
