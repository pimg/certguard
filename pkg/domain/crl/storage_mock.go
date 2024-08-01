package crl

import (
	"context"
	"crypto/x509"
)

// TODO create better mock repository that can be used for testing
type MockRepository struct{}

func (r *MockRepository) List(_ context.Context) ([]*CertificateRevocationList, error) {
	return []*CertificateRevocationList{}, nil
}

func (r *MockRepository) Save(_ context.Context, _ *CertificateRevocationList) (int64, error) {
	return 0, nil
}

func (r *MockRepository) Find(_ context.Context, _ string) (*CertificateRevocationList, error) {
	return &CertificateRevocationList{}, nil
}

func (r *MockRepository) SaveRevokedCertificates(_ context.Context, _ int64, _ []x509.RevocationListEntry) (int, error) {
	return 0, nil
}

func (r *MockRepository) FindRevokedCertificates(_ context.Context, _ int64) ([]*RevokedCertificate, error) {
	return []*RevokedCertificate{}, nil
}

func (r *MockRepository) Delete(_ context.Context, _ int64) error {
	return nil
}

func NewMockStorage() (*Storage, error) {
	return NewStorage(&MockRepository{}, "test")
}
