package client

import (
	"fmt"
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
		go common.OnceTimerFunc(func() {
			log.Warn("Connect to target server error, wait retry...")
			_ = s.Start()
		}).Run(10 * time.Second)
		return nil
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
	body := protocol.ClientMessageBody{
		Type:      Heartbeat,
		Timestamp: uint64(time.Now().UnixMilli()),
	}

	buffer := buf.NewBuffer(protocol.ClientMessageBodyLen)
	defer buffer.Release()

	_ = body.WriteTo(buffer)
	if err := tcpClient.Write(buffer); err != nil {
		log.Error("Send heartbeat packet error, ", err)
		tcpClient.Disconnect()
		reconnect(tcpClient)
		return
	}
	log.Debug("Send heartbeat packet success")
}

func handshake(ctx context.Context, conn net.Conn) error {
	timeout := time.Now().Add(HandshakeTimeout)
	buffer := buf.NewBuffer(256)
	ctx = context.WithTimeout(&ctx, HandshakeTimeout)

	defer func() {
		ctx.Cancel()
		buffer.Release()
	}()

	{
		authBody := protocol.ClientAuth{
			Name:     ClientCnf.Name,
			Internal: net.ParseFromStr(ClientCnf.Internal).Bytes(),
			Flag:     0,
			Key:      ClientCnf.Key,
		}
		header := protocol.NewAuthHeader(tiangong.VersionByte(), protocol.AuthClient)
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
		return errors.NewError("Handshake Timeout", ctx.Err())
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

	if err := conn.SetWriteDeadline(time.Time{}); err != nil {
		return errors.NewError("SetWriteDeadline error", err)
	}

	log.Info("Connect to target server [%s] seuccess", conn.RemoteAddr().String())
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
