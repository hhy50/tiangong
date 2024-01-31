package proxy_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/haiyanghan/tiangong/common/buf"
	"github.com/haiyanghan/tiangong/common/log"
	"github.com/haiyanghan/tiangong/common/net"
	"github.com/haiyanghan/tiangong/kernel/proxy"
)

var (
	ProxyHost = "127.0.0.1"
	ProxyPort = 8080
)

func TestHttpPorxy(t *testing.T) {
	log.InitLog()
	startTcpServer()
	if err := proxy.SetProxy(fmt.Sprintf("%s:%d", ProxyHost, ProxyPort), []string{"192.168.110.*"}); err != nil {
		t.Error(err)
		return
	}
	<-time.NewTimer(20 * time.Second).C
  	if err := proxy.ResetProxy(); err != nil {
		t.Error(err)
		return
	}
}

func startTcpServer() {
	ctx, _ := context.WithCancel(context.Background())
	tcpServer := net.NewTcpServer(ProxyHost, ProxyPort, ctx)
	connFunc := func(ctx context.Context, conn net.Conn) error {
		go func() {
			buffer := buf.NewRingBuffer()
			for {
				if _, err := buffer.Write(conn, buffer.Cap()); err != nil {
					log.Error("read error", err)
					conn.Close()
					return
				}
				bytes, _ := buf.ReadAll(buffer)
				log.Info("Rec %d byes, [%s]", len(bytes), bytes)
			}
		}()
		return nil
	}

	if err := tcpServer.ListenTCP(connFunc); err != nil {
		panic(err)
	}
}
