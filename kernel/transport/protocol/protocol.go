package protocol

import "tiangong/common/buf"

const (
	TCP int = iota
	UDP
	ICMP
)

type Segment []byte

type InBound interface {
	Decode([]byte) error
}

type OutBound interface {
	Encode() (buf.Buffer, error)
}
