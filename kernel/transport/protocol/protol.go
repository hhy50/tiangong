package protocol

import (
	"google.golang.org/protobuf/proto"
	"io"
)

func DecodeProtoMessage(reader io.Reader, length int, message proto.Message) error {
	bytes := make([]byte, length)
	if n, err := reader.Read(bytes); err != nil || n != length {
		return err
	}
	if err := proto.Unmarshal(bytes, message); err != nil {
		return err
	}
	return nil
}
