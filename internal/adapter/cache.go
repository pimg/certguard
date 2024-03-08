package adapter

type Cache interface {
	Write(name string, contents []byte) error
	Read(name string) ([]byte, error)
}

var GlobalCache Cache
