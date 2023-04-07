package writer

type Writer interface {
	Write(bytes []byte, filename string) error
}
