package protocol

import (
	"fmt"
	"io"
	"math"

	"github.com/haiyanghan/tiangong/common/buf"
	"github.com/haiyanghan/tiangong/common/errors"
	"google.golang.org/protobuf/proto"
)

type AuthType = byte
type AuthStatus = byte

const (
	AuthHeaderLen = 16

	// AuthType Define
	AuthSession AuthType = 1
	AuthClient  AuthType = 2
)

// AuthHeader byte length is AuthHeaderLen
type AuthHeader struct {
	Version  byte     // Client kernel version
	Type     AuthType // AuthType
	Reserved [13]byte // Reserved
	Len      byte     // DataLength

	body proto.Message
}

func NewAuthHeader(version byte, t AuthType) *AuthHeader {
	return &AuthHeader{
		Version:  version,
		Type:     t,
		Reserved: [13]byte{},
	}
}

func (h *AuthHeader) AppendBody(m proto.Message) *AuthHeader {
	h.body = m
	return h
}

func (h *AuthHeader) WriteTo(buffer buf.Buffer) error {
	if h.body == nil {
		return errors.NewError("Auth body is null", nil)
	}
	var body []byte
	{
		b, err := proto.Marshal(h.body)
		if err != nil {
			return err
		}
		if h.Len = byte(len(b)); int(h.Len) != len(b) {
			return errors.NewError(fmt.Sprintf("AuthBody Len too long, max limit: %d", math.MaxUint8), nil)
		}
		body = b
	}

	if err := buf.WriteByte(buffer, h.Version); err != nil {
		return err
	}
	if err := buf.WriteByte(buffer, h.Type); err != nil {
		return err
	}
	if err := buf.WriteBytes(buffer, h.Reserved[:]); err != nil {
		return err
	}
	if err := buf.WriteByte(buffer, h.Len); err != nil {
		return err
	}
	if err := buf.WriteBytes(buffer, body); err != nil {
		return err
	}
	return nil
}

func DecodeAuthHeader(buffer buf.Buffer, header *AuthHeader) error {
	if buffer.Len() < AuthHeaderLen {
		return io.EOF
	}

	header.Version, _ = buf.ReadByte(buffer)
	header.Type, _ = buf.ReadByte(buffer)
	{
		// Skip Reserved
		for i := range header.Reserved {
			header.Reserved[i], _ = buf.ReadByte(buffer)
		}
	}
	header.Len, _ = buf.ReadByte(buffer)
	return nil
}
