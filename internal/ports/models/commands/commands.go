package commands

import (
	domain_crl "github.com/pimg/certguard/pkg/domain/crl"
)

type Commands struct {
	storage *domain_crl.Storage
}

func NewCommands(storage *domain_crl.Storage) *Commands {
	return &Commands{
		storage: storage,
	}
}

func (c *Commands) CacheDir() string {
	return c.storage.CacheDir()
}

func (c *Commands) ImportDir() string {
	return c.storage.ImportDir()
}
