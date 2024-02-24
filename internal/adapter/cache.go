package adapter

type Cache interface {
	Write(name string, contents []byte) error
}

var GlobalCache Cache
