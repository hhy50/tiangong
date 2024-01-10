package net

// A four byte IP address
type IpAddress []byte

var (
	Internal IpAddress
	Local    IpAddress = []byte{127, 0, 0, 1}
)

func init() {
	// TODO
	// Get Internal ip From
}

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
