package main

import (
	"context"
	"fmt"
	"math"
	"runtime"
	"time"

	"github.com/haiyanghan/tiangong"
	"github.com/haiyanghan/tiangong/common"
	"github.com/haiyanghan/tiangong/common/buf"
	"github.com/haiyanghan/tiangong/common/errors"
	"github.com/haiyanghan/tiangong/common/log"
	"github.com/haiyanghan/tiangong/common/net"
	"github.com/haiyanghan/tiangong/transport/protocol"
)

type Resource struct {
	net.TcpClient
	num   int
	incre common.Incrementer
}

func (r *Resource) WriteToRemote(proto protocol.Protocol, bytes buf.Buffer) error {
	l := bytes.Len()
	if l > math.MaxUint16 {
		return errors.NewError("packet length exceeding maximum limit", nil)
	}
	h := protocol.PacketHeader{
		Len:      uint16(l),
		Rid:      uint32(r.incre.Next()),
		Protocol: byte(proto),
	}

	buffer := buf.NewBuffer(protocol.PacketHeaderLen + l)
	defer buffer.Release()

	// write packet header
	_ = h.WriteTo(buffer)
	// write body
	_, _ = buffer.Write(bytes, l)
	return r.Write(buffer)
}

func ConnSuccess(ctx context.Context, conn net.Conn) error {
	closeFunc := conn.Close
	// handshake
	if err := handshake(conn, Config.Token, net.ParseIp(Config.Subhost)); err != nil {
		_ = closeFunc()
		return err
	}
	log.Info("Connect to target server success [%s]", conn.RemoteAddr())
	return nil
}

func StartListener(ctx context.Context, conn net.Conn) error {
	go func() {
		client, _ := GetResourceWithTimeout(time.Hour)
		buffer := buf.NewRingBuffer()
		defer buffer.Release()
		defer conn.Close()
		defer PutResource(client)

		for {
			select {
			case <-ctx.Done():
				runtime.Goexit()
			default:
				if n, err := buffer.Write(conn, buffer.Cap()); err != nil || n <= 0 {
					log.Warn("Connect closed, %s", err.Error())
					runtime.Goexit()
				}

				bytes, _ := buf.ReadAll(buffer)
				if err := client.WriteToRemote(protocol.TCP, buf.Wrap(bytes)); err != nil {
					log.Error("Write to remote server error", err)
				}
				common.OnceTimerFunc(func() {
					runtime.Goexit()
				}).Run(2 * time.Second)
			}
		}
	}()

	return nil
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
		header := protocol.NewAuthHeader(tiangong.VersionByte(), protocol.AuthSession)
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
