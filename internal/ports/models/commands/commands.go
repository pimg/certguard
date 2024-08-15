package commands

import (
	domain_crl "github.com/pimg/certguard/pkg/domain/crl"
)

type Commands struct {
	Storage *domain_crl.Storage
}

func NewCommands(storage *domain_crl.Storage) *Commands {
	return &Commands{
		Storage: storage,
	}
}
