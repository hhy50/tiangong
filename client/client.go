package client

import (
	"strconv"
	"strings"
	"time"

	"github.com/haiyanghan/tiangong/common/context"

	"github.com/haiyanghan/tiangong"
	"github.com/haiyanghan/tiangong/common"
	"github.com/haiyanghan/tiangong/common/buf"
	"github.com/haiyanghan/tiangong/common/conf"
	"github.com/haiyanghan/tiangong/common/errors"
	"github.com/haiyanghan/tiangong/common/log"
	"github.com/haiyanghan/tiangong/common/net"
	"github.com/haiyanghan/tiangong/transport/protocol"
)

var (
	ConnTimeout      = 30 * time.Second
	HandshakeTimeout = ConnTimeout
	ClientCnf        Config
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
		// go common.OnceTimerFunc(func() {
		// 	log.Warn("Connect to target server error, wait retry...")
		// 	_ = s.Start()
		// }).Run(10 * time.Second)
		// return nil
	}
	go common.TimerFunc(func() {
		heartbeat(s.tcpClient)
	}).Run(time.Minute)
	return nil
}

func (s *clientImpl) Stop() {

}

func heartbeat(tcpClient net.TcpClient) {
	if !tcpClient.Connected() {
		reconnect(tcpClient)
		return
	}
	buffer := buf.NewBuffer(protocol.PacketHeaderLen)
	defer buffer.Release()

	heartbeatPacket := protocol.NewHeartbeatPacket()
	if err := protocol.EncodePacket(buffer, heartbeatPacket); err != nil {
		log.Error("send heartbeat packet error, ", err)
		return
	}
	if err := tcpClient.Write(buffer); err != nil {
		log.Error("Send heartbeat packet error, ", err)
		_ = tcpClient.Disconnect()
		reconnect(tcpClient)
		return
	}
	log.Debug("Send heartbeat packet success")
}

func handshake(ctx context.Context, conn net.Conn) error {
	timeout := time.Now().Add(HandshakeTimeout)
	ctx = context.WithTimeout(ctx, HandshakeTimeout)

	buffer := buf.NewBuffer(4096)
	defer func() {
		ctx.Cancel()
		buffer.Release()
	}()

	{
		packet, err := protocol.NewAuthRequestPacket(tiangong.VersionByte(), protocol.AuthClient, &protocol.ClientAuthBody{
			Name:     ClientCnf.Name,
			Internal: ClientCnf.Internal,
			Key:      ClientCnf.Key,
			Export:   ClientCnf.Export,
		})
		if err != nil {
			return err
		}
		if err := protocol.EncodePacket(buffer, packet); err != nil {
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
		return errors.NewError("Handshake Timeout", ctx.Err())
	default:
		// AuthResponsePacket
		header, err := protocol.DecodePacket(buffer, conn, timeout.Sub(time.Now()))
		if err != nil || !header.AuthSuccess() {
			return errors.NewError("handshake fail", err)
		}
	}

	if err := conn.SetWriteDeadline(time.Time{}); err != nil {
		return errors.NewError("SetWriteDeadline error", err)
	}

	log.Info("Connect to target server [%s] success", conn.RemoteAddr().String())
	log.Info("Handshake success")
	return nil
}

func reconnect(tcpClient net.TcpClient) {
	err := tcpClient.Connect(handshake)
	if err != nil {
		log.Warn("Reconnect to target server error, %s", err.Error())
		return
	}
	log.Info("Reconnect to target server success.")
}

// NewClient
func NewClient() (Client, error) {
	if err := conf.LoadConfig("", &ClientCnf); err != nil {
		return nil, err
	}

	if err := ClientCnf.Require(); err != nil {
		return nil, err
	}

	ctx := context.Empty()
	serverAddr := strings.Split(ClientCnf.Address, ":")
	serverPort, _ := strconv.Atoi(serverAddr[1])

	return &clientImpl{
		ctx:       ctx,
		tcpClient: net.NewTcpClient(serverAddr[0], serverPort, ctx),
	}, nil
}
