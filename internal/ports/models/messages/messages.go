package messages

import "crypto/x509"

type CRLResponseMsg struct {
	RevocationList *x509.RevocationList
}

type ErrorMsg struct {
	Err error
}
