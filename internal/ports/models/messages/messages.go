package messages

import (
	"crypto/x509"
	"net/url"

	"github.com/pimg/certguard/pkg/domain/crl"
)

type CRLResponseMsg struct {
	RevocationList *x509.RevocationList
	URL            *url.URL
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

type PemCertificateMsg struct {
	Certificate *x509.Certificate
}

type GetRevokedCertificateMsg struct {
	RevokedCertificate *crl.RevokedCertificate
	Found              bool
}
