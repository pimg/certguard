package crl

import (
	"context"
	"crypto/x509"
	"errors"
	"net/url"
	"time"
)

type CertificateRevocationList struct {
	ID         int64
	Name       string
	Signature  []byte
	ThisUpdate time.Time
	NextUpdate time.Time
	Raw        []byte
	URL        *url.URL
}

// TODO changeto map[int]RevocationReason
var RevocationReasons = map[int]string{
	0:  "unspecified",
	1:  "keyCompromise",
	2:  "cACompromise",
	3:  "affiliationChanged",
	4:  "superseded",
	5:  "cessationOfOperation",
	6:  "certificateHold",
	8:  "removeFromCRL",
	9:  "priviledgeWithdrawn",
	10: "aACompromise",
}

func FromCRL(crl *x509.RevocationList, URL *url.URL) (*CertificateRevocationList, error) {
	return &CertificateRevocationList{
		Name:       crl.Issuer.CommonName,
		Signature:  crl.Signature,
		ThisUpdate: crl.ThisUpdate,
		NextUpdate: crl.NextUpdate,
		URL:        URL,
		Raw:        crl.Raw,
	}, nil
}

func Process(ctx context.Context, URL *url.URL, crl *x509.RevocationList, store *Storage) error {
	parsed, err := FromCRL(crl, URL)
	if err != nil {
		return err
	}

	id, err := store.Repository.Save(ctx, parsed)
	if err != nil {
		return err
	}

	storedRevokedCertificates, err := store.Repository.SaveRevokedCertificates(ctx, id, crl.RevokedCertificateEntries)
	if err != nil {
		return err
	}

	if storedRevokedCertificates < len(crl.RevokedCertificateEntries) {
		return errors.New("no revoked certificates saved")
	}

	return nil
}
