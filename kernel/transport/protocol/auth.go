package protocol

import (
	"bytes"
	"strconv"
	"tiangong/common/buf"
	"tiangong/common/errors"
	"tiangong/common/net"
)

type Flag = byte

const (
	Hide         = 0x0001
	MAX_NAME_LEN = 32
)

type Auth struct {
	NameLen  byte          // Name Actual length
	Name     string        // ClientName
	Internal net.IpAddress // Client Custom InternalIp
	Flag     Flag          // Flag, 0: Hide,
}

func (a *Auth) Encode() (buf.Buffer, error) {
	name := bytes.NewBufferString(a.Name)
	namelen := name.Len()
	if namelen > MAX_NAME_LEN {
		return nil, errors.NewError("Name exceeds the maximum length limit, maximum length: "+strconv.Itoa(MAX_NAME_LEN), nil)
	}
	buffer := buf.NewBuffer(1 + namelen + 1 + 1)
	_ = buf.WriteByte(buffer, byte(namelen))
	_, _ = buf.WriteBytes(buffer, name)
	_, _ = buf.WriteBytes(buffer, bytes.NewBuffer(a.Internal[:]))
	_ = buf.WriteByte(buffer, a.Flag)
	return buffer, nil
}

func (a *Auth) Decode(bytes []byte) error {
	return nil
}
