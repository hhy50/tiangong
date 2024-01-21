package server

import (
	"context"
	"tiangong/common"
	"tiangong/common/net"
	"tiangong/kernel/transport/protocol"
	"tiangong/server/auth"
	"tiangong/server/client"
	"tiangong/server/internal"
	"tiangong/server/session"

	"github.com/google/uuid"
)

type Cli = *protocol.ClientAuth
type Session = *protocol.SessionAuth

type ListenFunc func()

func (l ListenFunc) Run() { l() }

func CloseConn(conn net.Conn) {
	_ = conn.Close()
}

func connHandler(ctx context.Context, conn net.Conn) error {
	_, user, err := auth.Authentication(ServerCnf.Key, conn)
	if err != nil {
		CloseConn(conn)
		return err
	}

	var runner common.Runnable
	switch user.(type) {
	case Cli:
		cli := user.(Cli)
		c := buildClient(conn, cli)
		_ = client.RegistClient(&c)
		runner = ListenFunc(func() {
		})
		break
	case Session:
		subHost := net.ValueOf((user.(Session)).SubHost)
		s := session.NewSession(subHost, (user.(Session)).Token, conn, ctx)
		if err = session.AddSession(&s); err != nil {
			return err
		}
		runner = ListenFunc(s.Work)
		break
	}

	go runner.Run()
	return nil
}

func buildClient(conn net.Conn, cli Cli) client.Client {
	getInternalIpFromReq := func() net.IpAddress {
		if len(cli.Internal) == 4 {
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
	return client.NewClient(internalIP, cli, conn)
}
