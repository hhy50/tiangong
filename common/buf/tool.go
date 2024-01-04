package buf

import (
	"io"
)

func WriteInt(buffer Buffer, data int) error {
	return nil
}

func WriteByte(buffer Buffer, data byte) error {
	return nil
}

func WriteBytes(buffer Buffer, reader io.Reader) (int, error) {
	return buffer.Write(reader)
}

func ReadInt(buffer Buffer) (int, error) {
	return 0, nil
}

func ReadByte(buffer Buffer) (byte, error) {
	return 0, nil
}
