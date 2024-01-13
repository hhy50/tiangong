package buf

type Reader interface {
	ReadFrom(buffer Buffer) error
}
