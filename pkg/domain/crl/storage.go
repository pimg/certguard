package crl

import (
	"context"
	"crypto/x509"
)

type Repository interface {
	Save(ctx context.Context, crl *CertificateRevocationList) (int64, error)
	Find(ctx context.Context, name string) (*CertificateRevocationList, error)
	List(ctx context.Context) ([]*CertificateRevocationList, error)
	Delete(ctx context.Context, id int64) error
	SaveRevokedCertificates(ctx context.Context, revocationListId int64, revokedCertificates []x509.RevocationListEntry) (int, error)
	FindRevokedCertificates(ctx context.Context, revocationListId int64) ([]*RevokedCertificate, error)
	FindRevokedCertificate(ctx context.Context, serialnumber string) (*RevokedCertificate, error)
}

type Storage struct {
	Repository Repository
	baseDir    string
	importDir  string
}

func NewStorage(repository Repository, baseDir, importDir string) (*Storage, error) {
	storage := &Storage{
		Repository: repository,
		baseDir:    baseDir,
		importDir:  importDir,
	}
	return storage, nil
}

func (s *Storage) CacheDir() string {
	return s.baseDir
}

func (s *Storage) ImportDir() string {
	return s.importDir
}
