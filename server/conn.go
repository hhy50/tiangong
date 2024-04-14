package server

import (
	"context"
	"reflect"

	"github.com/haiyanghan/tiangong/common"
	"github.com/haiyanghan/tiangong/common/net"
	"github.com/haiyanghan/tiangong/server/auth"
	"github.com/haiyanghan/tiangong/server/client"
	"github.com/haiyanghan/tiangong/server/internal"
	"github.com/haiyanghan/tiangong/server/session"
	"github.com/haiyanghan/tiangong/transport/protocol"

	"github.com/google/uuid"
)

var (
	NoAlloc = net.IpAddress{0, 0, 0, 0}
)

type Cli = *protocol.ClientAuth
type Session = *protocol.SessionAuth

type ListenFunc func()

func (l ListenFunc) Run() { l() }

func connHandler(ctx context.Context, conn net.Conn) error {
	_, user, err := auth.Authentication(ServerCnf.Key, conn)
	if err != nil {
		return err
	}

	var runner common.Runnable
	switch user.(type) {
	case Cli:
		cli := user.(Cli)
		c := buildClient(ctx, conn, cli)
		if err := client.RegistClient(&c); err != nil {
			return err
		}
		runner = ListenFunc(c.Keepalive)
	case Session:
		subHost := net.ParseFromBytes((user.(Session)).SubHost)
		s := session.NewSession(subHost, (user.(Session)).Token, conn, ctx)
		if err = session.AddSession(&s); err != nil {
			return err
		}
		runner = ListenFunc(s.Work)
	}

	go runner.Run()
	return nil
}

func buildClient(ctx context.Context, conn net.Conn, cli Cli) client.Client {
	getInternalIpFromReq := func() net.IpAddress {
		if len(cli.Internal) == 4 || reflect.DeepEqual(cli.Internal, NoAlloc) {
			i := cli.Internal
			return net.IpAddress{i[0], i[1], i[2], i[3]}
		}
		return internal.GeneraInternalIp()
	}

	internalIP := getInternalIpFromReq()
	if common.IsEmpty(cli.Name) {
		uid, _ := uuid.NewUUID()
		cli.Name = uid.String()
	}
	return client.NewClient(ctx, internalIP, cli, conn)
}
