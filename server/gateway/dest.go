package gateway

type Destination interface {
}

type DirectDestination struct {
	Host string
	Port Port
}

type Port [4]byte

func (p *Port) GetA() byte {
	return p[0]
}

func (p *Port) GetB() byte {
	return p[0]
}

func (p *Port) GetC() byte {
	return p[0]
}

func (p *Port) GetD() byte {
	return p[0]
}
