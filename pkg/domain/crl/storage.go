package crl

import (
	"context"
	"crypto/x509"
	"path/filepath"
)

type Repository interface {
	Save(ctx context.Context, crl *CertificateRevocationList) (int64, error)
	Find(ctx context.Context, name string) (*CertificateRevocationList, error)
	List(ctx context.Context) ([]*CertificateRevocationList, error)
	SaveRevokedCertificates(ctx context.Context, revocationListId int64, revokedCertificates []x509.RevocationListEntry) (int, error)
	FindRevokedCertificates(ctx context.Context, revocationListId int64) ([]*RevokedCertificate, error)
}

type Storage struct {
	Repository Repository
	baseDir    string
}

func NewStorage(repository Repository, baseDir string) (*Storage, error) {
	storage := &Storage{
		Repository: repository,
		baseDir:    baseDir,
	}
	GlobalStorage = storage
	return storage, nil
}

var GlobalStorage *Storage

func (s *Storage) CacheDir() string {
	return s.baseDir
}

func (s *Storage) ImportDir() string {
	return filepath.Join(s.baseDir, "/import")
}
