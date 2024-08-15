package db

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/pimg/certguard/pkg/domain/crl"
)

// Find all revoked certificates in CRL
func (s *LibSqlStorage) FindRevokedCertificates(ctx context.Context, revocationListID int64) ([]*crl.RevokedCertificate, error) {
	log.Printf("find revoked certificates in CRL: %d", revocationListID)
	dbRevCerts, err := s.Queries.GetRevokedCertificatesByRevocationList(ctx, revocationListID)
	if err != nil {
		return nil, err
	}

	revokedCertificates := make([]*crl.RevokedCertificate, len(dbRevCerts))

	for i, revokedCertificate := range dbRevCerts {

		revocationDate, ok := revokedCertificate.RevocationDate.(time.Time)
		if !ok {
			return nil, errors.New("invalid revocation date")
		}

		revokedCertificates[i] = &crl.RevokedCertificate{
			SerialNumber:     revokedCertificate.Serialnumber,
			RevocationReason: crl.RevocationReason(revokedCertificate.Reason),
			RevocationDate:   revocationDate,
			RevocationListID: revokedCertificate.RevocationList,
		}
	}

	return revokedCertificates, nil
}

func (s *LibSqlStorage) FindRevokedCertificate(ctx context.Context, serialnumber string) (*crl.RevokedCertificate, error) {
	log.Printf("find revoked certificate by serial number: %s", serialnumber)
	dbRevokedCertificate, err := s.Queries.GetRevokedCertificate(ctx, serialnumber)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	if dbRevokedCertificate.Serialnumber != serialnumber {
		return nil, errors.New("invalid serial number")
	}

	revocationDate, ok := dbRevokedCertificate.RevocationDate.(time.Time)
	if !ok {
		return nil, errors.New("invalid revocation date")
	}

	return &crl.RevokedCertificate{
		SerialNumber:     dbRevokedCertificate.Serialnumber,
		RevocationReason: crl.RevocationReason(dbRevokedCertificate.Reason),
		RevocationDate:   revocationDate,
		RevocationListID: dbRevokedCertificate.ID,
		RevokedBy:        dbRevokedCertificate.RevokedBy,
	}, nil
}
