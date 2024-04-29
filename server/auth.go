package server

import (
	"encoding/json"
	"time"

	"github.com/haiyanghan/tiangong/common"
	"github.com/haiyanghan/tiangong/common/buf"
	"github.com/haiyanghan/tiangong/common/context"
	"github.com/haiyanghan/tiangong/common/errors"
	"github.com/haiyanghan/tiangong/common/log"
	"github.com/haiyanghan/tiangong/common/net"
	"github.com/haiyanghan/tiangong/server/client"
	"github.com/haiyanghan/tiangong/server/session"
	"github.com/haiyanghan/tiangong/transport/protocol"
)

var (
	Timout   = 15 * time.Second
	complete = func(conn net.Conn, status protocol.AuthStatus) {
		buffer := buf.NewBuffer(protocol.PacketHeaderLen)
		defer buffer.Release()

		packet := protocol.NewAuthResponsePacket(status)
		if err := protocol.EncodePacket(buffer, packet); err != nil {
			log.Warn("write to auth response error, %+v", err)
			return
		}
		if err := conn.ReadFrom(buffer); err != nil {
			log.Warn("Write to auth response error, %+v", err)
			return
		}
	}
)

func AuthKey(packet *protocol.Packet, ctx context.Context) error {
	authBody := protocol.ClientAuthBody{}
	if err := json.Unmarshal(packet.Body, &authBody); err != nil {
		return err
	}

	conn := ctx.Value(net.ConnValName).(net.Conn)
	cm := ctx.Value(client.ManagerName).(*client.Manager)
	if common.IsNotEmpty(Config.Key) && Config.Key != authBody.Key {
		complete(conn, protocol.AuthFail)
		return errors.NewError("Auth fail, client key not match", nil)
	}
	complete(conn, protocol.AuthSuccess)

	// Add to manager
	newClient := client.NewClient(ctx, &authBody)
	if err := cm.RegisterClient(&newClient); err != nil {
		return err
	}
	go newClient.Keepalive()
	return nil
}

func AuthToken(packet *protocol.Packet, ctx context.Context) error {
	authBody := protocol.SessionAuthBody{}
	if err := json.Unmarshal(packet.Body, &authBody); err != nil {
		return err
	}

	conn := ctx.Value(net.ConnValName).(net.Conn)
	if err := VerificationToken(authBody.Token, ctx); err != nil {
		complete(conn, protocol.AuthFail)
		return err
	}
	complete(conn, protocol.AuthSuccess)

	subhost := net.ParseFromStr(authBody.SubHost)
	cm := ctx.Value(client.ManagerName).(*client.Manager)
	sm := ctx.Value(session.ManagerName).(*session.Manager)

	if dstClient := cm.GetClient(subhost); dstClient != nil {
		// Add to manager
		newSession := session.NewSession(ctx, authBody.Token, dstClient)
		sm.AddSession(subhost, newSession)
		go newSession.Work()
	}
	return nil
}
