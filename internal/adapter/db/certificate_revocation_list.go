package db

import (
	"context"
	"crypto/x509"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/pimg/certguard/internal/adapter/db/queries"
	"github.com/pimg/certguard/pkg/domain/crl"
)

// insert record
func (s *LibSqlStorage) Save(ctx context.Context, crl *crl.CertificateRevocationList) (int64, error) {
	params := queries.CreateCertificateRevocationListParams{
		Name:       crl.Name,
		Signature:  crl.Signature,
		ThisUpdate: crl.ThisUpdate,
		NextUpdate: sql.NullTime{
			Time:  crl.NextUpdate,
			Valid: true,
		},
		Raw: crl.Raw,
	}

	if crl.URL != nil {
		params.Url = sql.NullString{
			String: crl.URL.String(),
			Valid:  true,
		}
	}

	id, err := s.Queries.CreateCertificateRevocationList(ctx, params)
	if err != nil {
		log.Println("could not create certificate revocation list")
		return 0, err
	}

	fmt.Printf("crl with id: %d stored\n", id)
	return id, nil
}

// Find a Certificate Revocation List
func (s *LibSqlStorage) Find(ctx context.Context, name string) (*crl.CertificateRevocationList, error) {
	dbCrl, err := s.Queries.GetCertificateRevocationList(ctx, name)
	if err != nil {
		return nil, err
	}

	revocationList := &crl.CertificateRevocationList{
		ID:        dbCrl.ID,
		Name:      dbCrl.Name,
		Signature: dbCrl.Signature,
		Raw:       dbCrl.Raw,
	}

	nextUpdate, ok := dbCrl.NextUpdate.(time.Time)
	if !ok {
		return nil, errors.New("next_update not valid")
	}
	revocationList.NextUpdate = nextUpdate

	thisUpdate, ok := dbCrl.ThisUpdate.(time.Time)
	if !ok {
		return nil, errors.New("this_update not valid")
	}
	revocationList.ThisUpdate = thisUpdate

	return revocationList, nil
}

// List all Certificate Revocation Lists
func (s *LibSqlStorage) List(ctx context.Context) ([]*crl.CertificateRevocationList, error) {
	dbCrls, err := s.Queries.ListCertificateRevocationLists(ctx)
	if err != nil {
		return nil, err
	}

	cRLs := make([]*crl.CertificateRevocationList, len(dbCrls))

	for i, dbCrl := range dbCrls {
		thisUpdate, ok := dbCrl.ThisUpdate.(time.Time)
		if !ok {
			return nil, errors.New("invalid ThisUpdate")
		}

		nextUpdate, ok := dbCrl.NextUpdate.(time.Time)
		if !ok {
			return nil, errors.New("invalid NextUpdate")
		}

		url, err := url.Parse(dbCrl.Url.String)
		if err != nil {
			return nil, errors.Join(errors.New("invalid CRL URL"), err)
		}

		cRLs[i] = &crl.CertificateRevocationList{
			ID:         dbCrl.ID,
			Name:       dbCrl.Name,
			Signature:  dbCrl.Signature,
			ThisUpdate: thisUpdate,
			NextUpdate: nextUpdate,
			Raw:        dbCrl.Raw,
			URL:        url,
		}
	}

	return cRLs, nil
}

// save revoked certificates
// nolint: errcheck // checking err in defer results in panic
func (s *LibSqlStorage) SaveRevokedCertificates(ctx context.Context, revocationListId int64, revokedCertificates []x509.RevocationListEntry) (int, error) {
	tx, err := s.DB.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	qtx := s.Queries.WithTx(tx)

	rowsAffected := 0
	for _, revokedCertificateEntry := range revokedCertificates {
		reason, ok := crl.RevocationReasons[revokedCertificateEntry.ReasonCode]
		if !ok {
			return 0, errors.New("invalid ReasonCode on revoked certificate")
		}

		err := qtx.CreateRevokedCertificates(ctx, queries.CreateRevokedCertificatesParams{
			Serialnumber:   revokedCertificateEntry.SerialNumber.String(),
			RevocationDate: revokedCertificateEntry.RevocationTime,
			Reason:         reason.String(),
			RevocationList: revocationListId,
		})
		if err != nil {
			return 0, errors.Join(errors.New("could not save certificate revocation list entry"), err)
		}
		rowsAffected++
	}

	if rowsAffected != len(revokedCertificates) {
		return 0, errors.New("not all revoked certificates could be saved")
	}

	return rowsAffected, tx.Commit()
}

// delete certificate revocation list
func (s *LibSqlStorage) Delete(ctx context.Context, revocationListId int64) error {
	err := s.Queries.DeleteCertificateRevocationList(ctx, revocationListId)
	if err != nil {
		return errors.Join(errors.New("could not Delete CRL from storage "), err)
	}

	return nil
}
