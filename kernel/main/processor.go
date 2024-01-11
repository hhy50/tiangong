package main

import (
	"google.golang.org/protobuf/proto"
	"tiangong/common/log"
	"tiangong/common/net"
	"tiangong/kernel/transport/protocol"
	"time"
)

var (
	ConnTimeout = 30 * time.Second
)

type Processor struct {
	ServerHost net.IpAddress
	ServerPort net.Port
	Token      string
	SubHost    net.IpAddress
	client     net.TcpClient
}

func (p *Processor) ConnSuccess(conn net.Conn) {
	closeFunc := conn.Close
	// handshake
	if err := handshake(conn, p.Token, p.SubHost); err != nil {
		log.Error("handshake fail %+v", err)
		_ = closeFunc()
		return
	}
	log.Info("handshake success")

	//go func() {
	//	for {
	//		len, err := c.Read(buf)
	//		if len == 0 {
	//			fmt.Println("服务器已停止")
	//			return
	//		}
	//		if err != nil {
	//			fmt.Println("读取异常")
	//			return
	//		}
	//		msg := string(buf[:len-1])
	//		fmt.Println("收到服务端消息: ["+msg+"], 消息长度: ", len)
	//	}
	//}()
}

func NewProcessor() {

}

func handshake(conn net.Conn, token string, subHost net.IpAddress) error {
	a := protocol.SessionAuth{
		Token:   token,
		SubHost: subHost[:],
	}
	bts, err := proto.Marshal(&a)
	if err != nil {
		return err
	}
	if _, err := conn.Write(bts); err != nil {
		return err
	}
	return nil
}

func (p *Processor) Start() error {
	client := net.TcpClient{
		Host:    p.ServerHost,
		Port:    p.ServerPort,
		Timeout: ConnTimeout,
	}
	if err := client.Conn(p.ConnSuccess); err != nil {
		return err
	}
	return nil
}
