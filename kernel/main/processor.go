package main

import (
	"context"
	"fmt"
	"tiangong/common/errors"
	"tiangong/common/log"
	"tiangong/common/net"
	"tiangong/kernel"
	"tiangong/kernel/transport/protocol"
	"time"
)

var (
	ConnTimeout      = 30 * time.Second
	HandshakeTimeout = ConnTimeout
)

type Processor struct {
	ServerHost string
	ServerPort net.Port
	Token      string
	SubHost    net.IpAddress
	client     *net.TcpClient
}

func (p *Processor) ConnSuccess(conn net.Conn) error {
	closeFunc := conn.Close
	// handshake
	if err := handshake(conn, p.Token, p.SubHost); err != nil {
		_ = closeFunc()
		return err
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
	return nil
}

func init() {

}

func NewProcessor(host string, port int, token, sub string) Processor {
	return Processor{
		ServerHost: host,
		ServerPort: net.Port(port),
		SubHost:    net.ParseIp(sub),
		Token:      token,
	}
}

func handshake(conn net.Conn, token string, subHost net.IpAddress) error {
	timeout := time.Now().Add(HandshakeTimeout)
	ctx, cancel := context.WithTimeout(context.Background(), HandshakeTimeout)
	defer cancel()

	{
		authBody := protocol.SessionAuth{
			Token:   token,
			SubHost: subHost[:],
		}
		header := protocol.NewAuthHeader(kernel.VersionByte(), protocol.AuthSession)
		header.AppendBody(&authBody)
		bytes, err := header.ToBytes()
		if err != nil {
			return err
		}

		if err = conn.SetWriteDeadline(timeout); err != nil {
			return errors.NewError("SetWriteDeadline error", err)
		}
		if _, err = conn.Write(bytes); err != nil {
			return err
		}
	}
	select {
	case <-ctx.Done():
		return errors.NewError("Handshake Timeout", nil)
	default:
		if err := conn.SetReadDeadline(timeout); err != nil {
			return errors.NewError("SetReadDeadline error", err)
		}
		bytes := make([]byte, protocol.AuthResponseLen)
		if n, err := conn.Read(bytes); err != nil {
			return errors.NewError("", err)
		} else if n < protocol.AuthResponseLen {
			return errors.NewError(fmt.Sprintf("Auth response body too short, require %d bytes, Actual return %d bytes", protocol.AuthResponseLen, n), err)
		}

		response := protocol.AuthResponse{}
		if err := response.Unmarshal(bytes); err != nil || response.Status != protocol.AuthSuccess {
			return errors.NewError("handshake fail", err)
		}
	}

	return nil
}

func (p *Processor) Start() error {
	p.client = &net.TcpClient{
		Host:    p.ServerHost,
		Port:    p.ServerPort,
		Timeout: ConnTimeout,
	}
	if err := p.client.Conn(p.ConnSuccess); err != nil {
		return err
	}
	return nil
}
