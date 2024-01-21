package main

import (
	"context"
	"fmt"
	"tiangong/common"
	"tiangong/common/buf"
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
	ctx    context.Context
	client net.TcpClient
}

func ConnSuccess(ctx context.Context, conn net.Conn) error {
	closeFunc := conn.Close
	// handshake
	if err := handshake(conn, Token, net.ParseIp(SubHost)); err != nil {
		_ = closeFunc()
		return err
	}
	log.Info("connect success, target:[%s]", conn.RemoteAddr())

	p := common.GetProcess(ctx).(Processor)
	go func() {
		buffer := buf.NewRingBuffer()
		for {
			n, err := buffer.Write(buffer, buffer.Cap())
			if n == 0 || err != nil {
				_ = closeFunc()
				log.Error("connect closed...", err)
				p.RetryConnect()
				return
			}
		}
	}()
	return nil
}

func init() {

}

func NewProcessor() Processor {
	var p Processor

	ctx := common.SetProcess(context.Background(), p)
	p.ctx = ctx
	p.client = net.NewTcpClient(Server, Port, ctx)
	return p
}

func handshake(conn net.Conn, token string, subHost net.IpAddress) error {
	timeout := time.Now().Add(HandshakeTimeout)
	buffer := buf.NewBuffer(256)
	ctx, cancel := context.WithTimeout(context.Background(), HandshakeTimeout)
	defer cancel()
	defer buffer.Release()

	{
		authBody := protocol.SessionAuth{
			Token:   token,
			SubHost: subHost[:],
		}
		header := protocol.NewAuthHeader(kernel.VersionByte(), protocol.AuthSession)
		header.AppendBody(&authBody)
		if err := header.WriteTo(buffer); err != nil {
			return err
		}
		if err := conn.SetWriteDeadline(timeout); err != nil {
			return errors.NewError("SetWriteDeadline error", err)
		}

		if err := conn.ReadFrom(buffer); err != nil {
			return err
		}
		_ = buffer.Clear()
	}
	select {
	case <-ctx.Done():
		return errors.NewError("Handshake Timeout", nil)
	default:
		if err := conn.SetReadDeadline(timeout); err != nil {
			return errors.NewError("SetReadDeadline error", err)
		}
		if n, err := buffer.Write(conn, protocol.AuthResponseLen); err != nil {
			return errors.NewError("", err)
		} else if n < protocol.AuthResponseLen {
			return errors.NewError(fmt.Sprintf("Auth response body too short, require %d bytes, Actual return %d bytes", protocol.AuthResponseLen, n), err)
		}
		response := protocol.AuthResponse{}
		if err := response.ReadFrom(buffer); err != nil || response.Status != protocol.AuthSuccess {
			return errors.NewError("handshake fail", err)
		}
	}

	return nil
}

func (p *Processor) Start() error {
	if err := p.client.Connect(ConnSuccess); err != nil {
		// retry
		go p.RetryConnect()
		return nil
	}
	return nil
}

func (p *Processor) RetryConnect() {
	ticker := time.NewTicker(5 * time.Second)
	for {
		<-ticker.C
		if err := p.client.Connect(ConnSuccess); err != nil {
			log.Error("connect fail,", err)
		} else {
			ticker.Stop()
			break
		}
	}
}
