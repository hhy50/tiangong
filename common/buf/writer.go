package buf

type Writer interface {
	WriteTo(buffer Buffer) error
}
