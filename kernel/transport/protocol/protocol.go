package protocol

import "tiangong/common/buf"

type Segment []byte

type InBound interface {
	Decode([]byte) error
}

type OutBound interface {
	Encode() (buf.Buffer, error)
}
