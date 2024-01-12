package main

import (
	"context"
	"fmt"
	"google.golang.org/protobuf/proto"
	"tiangong/common/errors"
	"tiangong/common/log"
	"tiangong/common/net"
	"tiangong/kernel/transport/protocol"
	"time"
)

var (
	ConnTimeout      = 30 * time.Second
	HandshakeTimeout = ConnTimeout
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
	timeout := time.Now().Add(HandshakeTimeout)
	ctx, cancel := context.WithTimeout(context.Background(), HandshakeTimeout)
	defer cancel()

	a := protocol.SessionAuth{
		Token:   token,
		SubHost: subHost[:],
	}
	bts, err := proto.Marshal(&a)
	if err != nil {
		return err
	}
	if err := conn.SetWriteDeadline(timeout); err != nil {
		return errors.NewError("SetWriteDeadline error", err)
	}
	if _, err := conn.Write(bts); err != nil {
		return err
	}
	select {
	case <-ctx.Done():
		return errors.NewError("Handshake Timeout", nil)
	default:
		if err := conn.SetReadDeadline(timeout); err != nil {
			return errors.NewError("SetReadDeadline error", err)
		}

		bytes := make([]byte, protocol.AuthResponseLen)
		if n, err := conn.Read(bytes); err != nil || n < protocol.AuthResponseLen {
			return errors.NewError(fmt.Sprintf("Auth response body too short, require %d bytes, Actual return %d bytes", protocol.AuthResponseLen, n), err)
		}
		response, err := protocol.DecodeAuthResponse(bytes)
		if err != nil || response.Status != protocol.AuthSuccess {
			return errors.NewError("handshake fail", err)
		}
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
