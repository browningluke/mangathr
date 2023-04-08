package writer

import "fmt"

const partPrefix = ".part"

func getPartPath(path string) string {
	return fmt.Sprintf("%s%s", path, partPrefix)
}

type Writer interface {
	Write(bytes []byte, filename string) error
	MarkComplete() error
	Close() error
}
