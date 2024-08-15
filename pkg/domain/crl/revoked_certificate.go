package crl

import "time"

type RevokedCertificate struct {
	SerialNumber     string
	RevocationReason RevocationReason
	RevocationDate   time.Time
	RevocationListID int64
	RevokedBy        string
}

type RevocationReason string

func (r RevocationReason) String() string {
	return string(r)
}

const (
	RevocationReasonUnspecified          RevocationReason = "unspecified"
	RevocationReasonKeyCompromise        RevocationReason = "keyCompromise"
	RevocationReasonCACompromise         RevocationReason = "cACompromise"
	RevocationReasonAffiliationChanged   RevocationReason = "affiliationChanged"
	RevocationReasonSuperseded           RevocationReason = "superseded"
	RevocationReasonCessationOfOperation RevocationReason = "cessationOfOperation"
	RevocationReasonCertificateHold      RevocationReason = "certificateHold"
	RevocationReasonRemoveFromCRL        RevocationReason = "removeFromCRL"
	RevocationReasonPriviledgeWithdrawn  RevocationReason = "priviledgeWithdrawn"
	RevocationReasonAACompromise         RevocationReason = "aACompromise"
)
