package main

import (
	"context"
	"fmt"
	"math"
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
	ctx         context.Context
	client      net.TcpClient
	incrementer common.Incrementer
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
		defer buffer.Release()
		for {
			n, err := buffer.Write(buffer, buffer.Cap())
			if n == 0 || err != nil {
				_ = closeFunc()
				log.Error("connect closed...", err)
				go common.Retry(func() error {
					return p.client.Connect(ConnSuccess)
				}).Run(3*time.Second, -1)
				return
			}
		}
	}()
	return nil
}

func init() {

}

func NewProcessor() Processor {
	p := Processor{
		incrementer: common.Incrementer{Range: common.Range{0, math.MaxUint32}},
	}

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
		log.Error("connect tiangong-server error retry..., ", err)
		go common.Retry(func() error {
			return p.client.Connect(ConnSuccess)
		}).Run(3*time.Second, -1)
		return nil
	}
	return nil
}

func (p *Processor) WriteToRemote(proto byte, bytes buf.Buffer) error {
	l := bytes.Len()
	if l > math.MaxUint16 {
		return errors.NewError("packet length exceeding maximum limit", nil)
	}
	h := protocol.PacketHeader{
		Len:      uint16(l),
		Rid:      uint32(p.incrementer.Next()),
		Protocol: proto,
	}

	buffer := buf.NewBuffer(protocol.PacketHeaderLen + l)
	defer buffer.Release()

	_ = h.WriteTo(buffer)
	// write body
	_, _ = buffer.Write(bytes, l)
	return p.client.Write(buffer)
}
