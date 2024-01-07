package net

// A four byte IP address
type IpAddress []byte

var (
	Internal IpAddress = []byte{0x00, 0x00, 0x00, 0x00}
)

// Port Range 1024~65535
type Port uint16

func (p IpAddress) GetA() byte {
	return p[0]
}

func (p IpAddress) GetB() byte {
	return p[1]
}

func (p IpAddress) GetC() byte {
	return p[2]
}

func (p IpAddress) GetD() byte {
	return p[3]
}
