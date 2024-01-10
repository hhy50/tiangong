package buf

import (
	"tiangong/common"
)

func WriteByte(buffer Buffer, data byte) error {
	// TODO
	return nil
}

func WriteInt(buffer Buffer, data int) error {
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
