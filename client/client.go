package client

import (
	"context"
	"fmt"
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
)

type Client struct {
	Cnf       Config
	tcpClient net.TcpClient
}

func (s *Client) Start() error {
	if err := s.tcpClient.Connect(s.handshake); err != nil {
		return err
	}
	return nil
}
func (s *Client) Stop() error {
	s.tcpClient.Disconnect()
	return nil
}

func (c *Client) handshake(conn net.Conn) error {
	timeout := time.Now().Add(HandshakeTimeout)
	buffer := buf.NewBuffer(256)
	ctx, cancel := context.WithTimeout(context.Background(), HandshakeTimeout)

	defer cancel()
	defer buffer.Release()

	{
		authBody := protocol.ClientAuth{
			Name:     c.Cnf.Name,
			Internal: net.ParseIp(c.Cnf.Internal).Bytes(),
			Flag:     0,
			Key:      c.Cnf.Key,
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
func NewClient(cp string) (*Client, error) {
	c := Config{}
	if err := conf.LoadConfig(cp, &c, defaultValue); err != nil {
		return nil, err
	}

	if err := c.Require(); err != nil {
		return nil, err
	}

	tc := net.TcpClient{
		Host:    c.ServerHost,
		Port:    c.ServerPort,
		Timeout: 30 * time.Second,
	}
	return &Client{c, tc}, nil
}
