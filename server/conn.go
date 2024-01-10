package server

import (
	"github.com/google/uuid"
	"tiangong/common"
	"tiangong/common/net"
	tgNet "tiangong/common/net"
	"tiangong/kernel/transport/protocol"
	"tiangong/server/auth"
	"tiangong/server/client"
	"tiangong/server/internal"
	"tiangong/server/session"
)

type Cli = *protocol.ClientAuth
type Session = *protocol.SessionAuth

type ListenFunc func()

func (l ListenFunc) Run() { l() }

func connHandler(conn net.Conn) {
	close := func() {
		_ = conn.Close()
	}

	user, err := auth.Authentication(conn)
	if err != nil {
		close()
	}

	var runner common.Runnable
	switch user.(type) {
	case Cli:
		cli := user.(Cli)
		c := buildClient(conn, cli)
		_ = client.AddClient(&c)
		runner = ListenFunc(func() {

		})
		break
	case Session:
		ses := user.(Session)
		s := buildSession(conn, ses)
		_ = session.AddSession(&s)
		runner = ListenFunc(s.Work)
		break
	}

	go runner.Run()
}

func buildSession(conn net.Conn, ses Session) session.Session {
	if len(ses.SubHost) == 0 {
		ses.SubHost = ses.MainHost
	}
	return session.NewSession(ses.MainHost, ses.SubHost, ses.Token, conn)
}

func buildClient(conn net.Conn, cli Cli) client.Client {
	getInternalIpFromReq := func() tgNet.IpAddress {
		if len(cli.Internal) == 4 {
			return cli.Internal[0:4]
		}
		return nil
	}

	internalIp := getInternalIpFromReq()
	if internalIp != nil {
		cli.Internal = internal.GeneraInternalIp()
	}
	if common.IsEmpty(cli.Name) {
		uid, _ := uuid.NewUUID()
		cli.Name = uid.String()
	}
	return client.NewClient(cli, conn)
}
