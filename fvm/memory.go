package fvm

type Memory interface {
	Read(address uint16) byte
	Write(address uint16, value byte)
}

type Ram []byte

func (r Ram) Read(address uint16) byte {
	return r[address]
}

func (r Ram) Write(address uint16, value byte) {
	r[address] = value
}
