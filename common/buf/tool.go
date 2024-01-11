package buf

import (
	"bytes"
	"tiangong/common"
)

func WriteByte(buffer Buffer, b byte) error {
	return WriteBytes(buffer, []byte{b})
}

func WriteBytes(buffer Buffer, b []byte) error {
	if _, err := buffer.Write(bytes.NewBuffer(b), len(b)); err != nil {
		return err
	}
	return nil
}

func WriteInt(buffer Buffer, i int) error {
	// TODO
	return nil
}

func ReadByte(buffer Buffer) (byte, error) {
	one := [common.One]byte{}
	if n, err := buffer.Read(one[:]); err != nil || n != common.One {
		return 0, err
	}
	return one[0], nil
}

func ReadInt(buffer Buffer) (int, error) {
	// TODO
	return 0, nil
}
