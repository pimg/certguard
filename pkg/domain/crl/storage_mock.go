package crl

import (
	"context"
	"crypto/x509"
)

// TODO create better mock repository that can be used for testing
type MockRepository struct {
	CRLs                      map[int64]*CertificateRevocationList
	RevokedCertificateEntries map[int64][]x509.RevocationListEntry
}

func (r *MockRepository) FindRevokedCertificate(_ context.Context, _ string) (*RevokedCertificate, error) {
	return nil, nil
}

func (r *MockRepository) List(_ context.Context) ([]*CertificateRevocationList, error) {
	return []*CertificateRevocationList{}, nil
}

func (r *MockRepository) Save(_ context.Context, crl *CertificateRevocationList) (int64, error) {
	r.CRLs[crl.ID] = crl
	return crl.ID, nil
}

func (r *MockRepository) Find(_ context.Context, _ string) (*CertificateRevocationList, error) {
	return &CertificateRevocationList{}, nil
}

func (r *MockRepository) SaveRevokedCertificates(_ context.Context, crlID int64, entries []x509.RevocationListEntry) (int, error) {
	CRLentries, ok := r.RevokedCertificateEntries[crlID]
	if !ok {
		CRLentries = entries
	} else {
		CRLentries = append(CRLentries, entries...)
	}
	return len(CRLentries), nil
}

func (r *MockRepository) FindRevokedCertificates(_ context.Context, CRLID int64) ([]*RevokedCertificate, error) {
	CRL, ok := r.RevokedCertificateEntries[CRLID]
	revokedCertifcates := make([]*RevokedCertificate, 0)
	if !ok {
		return revokedCertifcates, nil
	}

	for _, entry := range CRL {
		revokedCertifcates = append(revokedCertifcates, &RevokedCertificate{
			SerialNumber:   entry.SerialNumber.String(),
			RevocationDate: entry.RevocationTime,
		})
	}
	return revokedCertifcates, nil
}

func (r *MockRepository) Delete(_ context.Context, _ int64) error {
	return nil
}

func NewMockStorage() (*Storage, error) {
	CRLs := make(map[int64]*CertificateRevocationList)
	RevokedCertificateEntries := make(map[int64][]x509.RevocationListEntry)
	return NewStorage(&MockRepository{
		CRLs:                      CRLs,
		RevokedCertificateEntries: RevokedCertificateEntries,
	}, "test")
}
