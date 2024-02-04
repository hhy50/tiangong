package main

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/haiyanghan/tiangong"
	"github.com/haiyanghan/tiangong/common/buf"
	"github.com/haiyanghan/tiangong/common/errors"
	"github.com/haiyanghan/tiangong/common/log"
	"github.com/haiyanghan/tiangong/common/net"
	"github.com/haiyanghan/tiangong/transport/protocol"
)

type Worker struct {
	conn   net.Conn
	client *Resource
	ctx    context.Context
}

func (w *Worker) Run() {
	buffer := buf.NewRingBuffer()
	defer func() {
		buffer.Release()
		w.conn.Close()
		if w.client != nil {
			PutResource(w.client)
		}
	}()

	for {
		select {
		case <-w.ctx.Done():
			return
		default:
			if n, err := buffer.Write(w.conn, buffer.Cap()); err != nil || n <= 0 {
				log.Warn("Connect closed, %s", err.Error())
				return
			}
			if w.client == nil {
				c, err := GetResourceWithTimeout(10 * time.Second)
				if err != nil {
					return
				}
				w.client = c
			}
			if err := w.client.WriteToRemote(protocol.TCP, buffer); err != nil {
				log.Error("Write to remote server error", err)
			}
		}
	}
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
	w := Worker{
		ctx:  ctx,
		conn: conn,
	}
	go w.Run()
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
