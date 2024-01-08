package client

import (
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
	"net"
	"tiangong/client"
	"tiangong/common"
	tgNet "tiangong/common/net"
	"tiangong/kernel/transport/protocol"
	"tiangong/server/auth"
	"tiangong/server/internal"
)

type Client struct {
	Name string
	Host tgNet.IpAddress

	header *protocol.AuthHeader
	conn   net.Conn
}

func NewClient(name string, host tgNet.IpAddress, header *protocol.AuthHeader, conn net.Conn) Client {
	return Client{
		Name:   name,
		Host:   host,
		header: header,
		conn:   conn,
	}
}

func ConnHandler(conn net.Conn) {
	close := func() {
		_ = conn.Close()
	}

	user, err := auth.Authentication(conn)
	if err != nil {
		close()
	}

	c := buildClient(conn, user)
	_ = client.AddClient(&c)

	switch user.(type) {
	//case client.Client:
	//	break
	//case session.Session:
	//	break
	}
}

func buildClient(conn net.Conn, user *proto.Message) client.Client {
	getInternalIpFromReq := func() tgNet.IpAddress {
		if len(auth.Internal) == 4 {
			return auth.Internal[0:4]
		}
		return nil
	}

	internalIp := getInternalIpFromReq()
	if internalIp != nil {
		internalIp = internal.GeneraInternalIp()
	}
	clientName := auth.Name
	if common.IsEmpty(clientName) {
		uid, _ := uuid.NewUUID()
		clientName = uid.String()
	}
	return client.NewClient(clientName, internalIp, conn)
}
