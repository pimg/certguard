package db

import (
	"context"
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
