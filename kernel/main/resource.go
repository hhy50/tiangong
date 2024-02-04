package main

import (
	"math"

	"github.com/haiyanghan/tiangong/common"
	"github.com/haiyanghan/tiangong/common/buf"
	"github.com/haiyanghan/tiangong/common/errors"
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
