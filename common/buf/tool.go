package buf

import (
	"encoding/binary"

	"github.com/haiyanghan/tiangong/common"
)

func WriteByte(buffer Buffer, b byte) error {
	return WriteBytes(buffer, []byte{b})
}

func WriteBytes(buffer Buffer, b []byte) error {
	if _, err := buffer.Write(Wrap(b), len(b)); err != nil {
		return err
	}
	return nil
}

func ReadByte(buffer Buffer) (byte, error) {
	return ReadUint8(buffer)
}

func ReadUint8(buffer Buffer) (uint8, error) {
	one := [common.One]uint8{}
	if n, err := buffer.Read(one[:]); err != nil || n != common.One {
		return 0, err
	}
	return one[0], nil
}

func ReadUint16(buffer Buffer) (uint16, error) {
	bytes := [common.Two]byte{}
	if n, err := buffer.Read(bytes[:]); err != nil || n != common.Two {
		return 0, err
	}
	return binary.BigEndian.Uint16(bytes[:]), nil
}

func ReadUint32(buffer Buffer) (uint32, error) {
	bytes := [common.Four]byte{}
	if n, err := buffer.Read(bytes[:]); err != nil || n != common.Four {
		return 0, err
	}
	return binary.BigEndian.Uint32(bytes[:]), nil
}

func ReadUint64(buffer Buffer) (uint64, error) {
	bytes := [common.Eight]byte{}
	if n, err := buffer.Read(bytes[:]); err != nil || n != common.Eight {
		return 0, err
	}
	return binary.BigEndian.Uint64(bytes[:]), nil
}

func ReadAll(buffer Buffer) ([]byte, error) {
	bytes := make([]byte, buffer.Len())
	if _, err := buffer.Read(bytes); err != nil {
		return nil, err
	}
	return bytes, nil
}
