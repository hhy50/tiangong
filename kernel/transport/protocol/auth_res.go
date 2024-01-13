package protocol

import (
	"strconv"
	"tiangong/common"
	"tiangong/common/errors"
	"time"
)

// AuthResponse byte length is AuthResponseLen
type AuthResponse struct {
	Status    AuthStatus
	Reserved  [7]byte
	Timestamp int64
}

func (r *AuthResponse) Marshal() ([]byte, error) {
	bytes := make([]byte, AuthResponseLen)
	bytes[0] = r.Status
	copy(bytes[1:8], r.Reserved[:])
	copy(bytes[8:], common.Uint64ToBytes(uint64(r.Timestamp)))
	return bytes, nil
}

func (r *AuthResponse) Unmarshal(bytes []byte) error {
	if len(bytes) < AuthResponseLen {
		return errors.NewError("bytes len too short, minnum is "+strconv.Itoa(AuthResponseLen)+"bytes", nil)
	}

	r.Status = bytes[0]
	r.Reserved = [7]byte{}
	r.Timestamp = int64(common.Uint64(bytes[8:]))
	return nil
}

func NewAuthResponse(status AuthStatus) *AuthResponse {
	return &AuthResponse{
		Status:    status,
		Reserved:  [7]byte{},
		Timestamp: time.Now().UnixMilli(),
	}
}
