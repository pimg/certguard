package adapter

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type FileCache struct {
	cacheDir string
}

func NewFileCache() error {
	cacheDir, err := determineCacheDir()
	if err != nil {
		return err
	}

	err = os.MkdirAll(cacheDir, 0o777)
	if err != nil {
		return err
	}

	GlobalCache = &FileCache{cacheDir: cacheDir}

	return nil
}

func (f *FileCache) Write(filename string, fileContent []byte) error {
	err := os.WriteFile(f.cacheDir+filename, fileContent, 0o777)
	if err != nil {
		return errors.Join(err, fmt.Errorf("failed to write filename to cache: %s", filename))
	}

	return nil
}

func (f *FileCache) Dir() string {
	return f.cacheDir
}

func determineCacheDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", errors.New("could not create file path to User home dir, Cache will not be enabled")
	}

	return filepath.Join(homeDir, ".cache", "certguard"), nil
}
