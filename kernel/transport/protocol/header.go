package protocol

type RequestHeader struct {
	magicHeader [2]byte
	version     [1]byte
	protocol    [1]byte
	uid         [8]byte
	len         [2]byte
	seg         Segment
	checksum	[2]byte
}
