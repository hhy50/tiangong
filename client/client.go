package client

import (
	"context"
	"fmt"
	"runtime"
	"tiangong/common"
	"tiangong/common/buf"
	"tiangong/common/conf"
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
	ClientCnf        Config

	// ClientMsgType
	Heartbeat byte = 1
)

type Client interface {
	Start() error
	Stop()
}

type clientImpl struct {
	ctx       context.Context
	tcpClient net.TcpClient
}

func (s *clientImpl) Start() error {
	if err := s.tcpClient.Connect(handshake); err != nil {
		return err
	}
	go common.TimerFunc(func() {
		body := protocol.ClientMessageBody{
			Type:      Heartbeat,
			Timestamp: uint64(time.Now().UnixMilli()),
		}

		buffer := buf.NewBuffer(protocol.ClientMessageBodyLen)
		defer buffer.Release()

		_ = body.WriteTo(buffer)
		if err := s.tcpClient.Write(buffer); err != nil {
			log.Error("send heartbeat packet error", err)
			runtime.Goexit()
		}
		log.Debug("send heartbeat packet success")
	}).Run(time.Minute)
	return nil
}
func (s *clientImpl) Stop() {
	cancel := s.ctx.Value(common.CancelFuncKey).(context.CancelFunc)
	cancel()
}

func handshake(ctx context.Context, conn net.Conn) error {
	timeout := time.Now().Add(HandshakeTimeout)
	buffer := buf.NewBuffer(256)
	ctx, cancel := context.WithTimeout(context.Background(), HandshakeTimeout)

	defer cancel()
	defer buffer.Release()

	{
		authBody := protocol.ClientAuth{
			Name:     ClientCnf.Name,
			Internal: net.ParseIp(ClientCnf.Internal).Bytes(),
			Flag:     0,
			Key:      ClientCnf.Key,
		}
		header := protocol.NewAuthHeader(kernel.VersionByte(), protocol.AuthClient)
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

	log.Info("handshake success")
	return nil
}

// NewClient by specify a config file
func NewClient(cp string) (Client, error) {
	if err := conf.LoadConfig(cp, &ClientCnf, defaultValue); err != nil {
		return nil, err
	}

	if err := ClientCnf.Require(); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	ctx = context.WithValue(ctx, common.CancelFuncKey, cancel)

	return &clientImpl{
		ctx:       ctx,
		tcpClient: net.NewTcpClient(ClientCnf.ServerHost, ClientCnf.ServerPort, ctx),
	}, nil
}
