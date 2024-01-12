package protocol

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	"io"
	"strconv"
	"tiangong/common"
	"tiangong/common/buf"
	"tiangong/common/errors"
	"time"
	"unsafe"
)

type AuthType = byte
type AuthStatus = byte

const (
	AuthHeaderLen   = int(unsafe.Sizeof((*AuthHeader)(nil)))
	AuthResponseLen = int(unsafe.Sizeof((*AuthResponse)(nil)))

	AuthFail    AuthStatus = 0
	AuthSuccess AuthStatus = 1
)

// AuthHeader byte length is 16
type AuthHeader struct {
	Version  byte     // Client kernel version
	Type     AuthType // AuthType
	Reserved [13]byte // Reserved
	Len      byte     // DataLength
}

// AuthResponse byte length is 16
type AuthResponse struct {
	Status    AuthStatus
	Reserved  [7]byte
	Timestamp int64
}

func (r *AuthResponse) Marshal() ([]byte, error) {
	return nil, nil
}

func NewAuthResponse(status AuthStatus) *AuthResponse {
	return &AuthResponse{
		Status:    status,
		Reserved:  [7]byte{},
		Timestamp: time.Now().UnixMilli(),
	}
}

func DecodeAuthHeader(reader io.Reader) (*AuthHeader, error) {
	bytes := [AuthHeaderLen]byte{}
	if n, err := reader.Read(bytes[:]); err != nil || n != AuthHeaderLen {
		return nil, errors.NewError(fmt.Sprintf("Auth fial, expect read %d bytes, actuality read %d bytes", AuthHeaderLen, n), err)
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

func DecodeClientAuthBody(reader io.Reader, bl byte) (proto.Message, error) {
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

func DecodeAuthResponse(bytes []byte) (*AuthResponse, error) {
	if len(bytes) < AuthResponseLen {
		return nil, errors.NewError("bytes len too short, minnum is "+strconv.Itoa(AuthResponseLen)+"bytes", nil)
	}

	reserved := [7]byte{}
	return &AuthResponse{
		Status:    bytes[0],
		Reserved:  reserved,
		Timestamp: int64(common.Uint64(bytes[8:17])),
	}, nil
}
