package messages

import (
	"crypto/x509"

	"github.com/pimg/certguard/pkg/domain/crl"
)

type CRLResponseMsg struct {
	RevocationList *x509.RevocationList
}

type ErrorMsg struct {
	Err error
}

type ListCRLsResponseMsg struct {
	CRLs []*crl.CertificateRevocationList
}

type RevokedCertificatesMsg struct {
	RevokedCertificates []x509.RevocationListEntry
}

type CRLDeleteConfirmationMsg struct {
	DeletionSuccessful bool
}
