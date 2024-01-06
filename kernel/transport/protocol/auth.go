package protocol

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	"io"
	"tiangong/common/buf"
	"tiangong/common/errors"
)

type AuthType = byte

const (
	authHeaderLen = 16
)

type AuthHeader struct {
	Version  byte     // Client kernel version
	Type     AuthType // AuthType
	Reserved [13]byte // Reserved
	Len      byte     // DataLength
}

func DecodeAuthHeader(reader io.Reader) (*AuthHeader, error) {
	bytes := [authHeaderLen]byte{}
	if n, err := reader.Read(bytes[:]); err != nil || n != authHeaderLen {
		return nil, errors.NewError(fmt.Sprintf("Auth fial, expect read %d bytes, actuality read %d bytes", authHeaderLen, n), err)
	}
	buffer := buf.Wrap(bytes[:])
	defer buffer.Release()

	header := AuthHeader{}
	header.Version, _ = buf.ReadByte(buffer)
	header.Type, _ = buf.ReadByte(buffer)
	{
		// Skip Reserved
		for i := range header.Reserved {
			header.Reserved[i], _ = buf.ReadByte(buffer)
		}
	}
	header.Len, _ = buf.ReadByte(buffer)
	return &header, nil
}

func DecodeClientAuthBody(reader io.Reader, bl byte) (*ClientAuth, error) {
	bytes := make([]byte, bl)
	if n, err := reader.Read(bytes); err != nil || n != int(bl) {
		return nil, err
	}
	auth := ClientAuth{}
	if err := proto.Unmarshal(bytes, &auth); err != nil {
		return nil, err
	}
	return &auth, nil
}
