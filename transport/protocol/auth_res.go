package protocol

import (
	"github.com/haiyanghan/tiangong/common"
	"github.com/haiyanghan/tiangong/common/buf"
	"github.com/haiyanghan/tiangong/common/errors"
	"strconv"
	"time"
)

// AuthResponse byte length is AuthResponseLen
type AuthResponse struct {
	Status    AuthStatus
	Reserved  [7]byte
	Timestamp uint64
}

func (r *AuthResponse) WriteTo(buffer buf.Buffer) error {
	if buffer.Cap() < AuthResponseLen {
		return errors.NewError("write bytes len too short, minnum is "+strconv.Itoa(AuthResponseLen)+"bytes", nil)
	}
	if err := buf.WriteByte(buffer, r.Status); err != nil {
		return err
	}
	if err := buf.WriteBytes(buffer, r.Reserved[:]); err != nil {
		return err
	}
	if err := buf.WriteBytes(buffer, common.Uint64ToBytes(r.Timestamp)); err != nil {
		return err
	}
	return nil
}

func (r *AuthResponse) ReadFrom(buffer buf.Buffer) error {
	if buffer.Len() < AuthResponseLen {
		return errors.NewError("read bytes len too short, minnum is "+strconv.Itoa(AuthResponseLen)+"bytes", nil)
	}

	r.Status, _ = buf.ReadByte(buffer)
	buffer.Read(r.Reserved[:]) // Skip Reserved
	r.Timestamp, _ = buf.ReadUint64(buffer)
	return nil
}

func NewAuthResponse(status AuthStatus) *AuthResponse {
	return &AuthResponse{
		Status:    status,
		Reserved:  [7]byte{},
		Timestamp: uint64(time.Now().UnixMilli()),
	}
}
