package commands

import (
	"context"
	"crypto/x509"
	"crypto/x509/pkix"
	"errors"
	"fmt"
	"log/slog"
	"math/big"
	"strconv"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pimg/certguard/internal/ports/models/messages"
	"github.com/pimg/certguard/pkg/domain/crl"
)

func GetRevokedCertificates(aID, aCN, aThisUpdate, aNextUpdate string) tea.Cmd {
	slog.Debug(fmt.Sprintf("getting revoked certificate from storage with ID: %s", aID))
	ctx := context.Background()
	return func() tea.Msg {
		ID, err := strconv.ParseInt(aID, 10, 64)
		if err != nil {
			slog.Debug(fmt.Sprintf("%s is not a valid ID, %v", aID, err))
			return messages.ErrorMsg{
				Err: errors.Join(errors.New("could not parse ID"), err),
			}
		}

		thisUpdate, err := time.Parse(time.DateOnly, aThisUpdate)
		if err != nil {
			slog.Debug(fmt.Sprintf("%s is not a valid time, %v", aThisUpdate, err))
			return messages.ErrorMsg{
				Err: errors.Join(errors.New("could not parse time of ThisUpdate"), err),
			}
		}

		nextUpdate, err := time.Parse(time.DateOnly, aNextUpdate)
		if err != nil {
			slog.Debug(fmt.Sprintf("%s is not a valid time, %v", aNextUpdate, err))
			return messages.ErrorMsg{
				Err: errors.Join(errors.New("could not parse time of NextUpdate"), err),
			}
		}

		certificates, err := crl.GlobalStorage.Repository.FindRevokedCertificates(ctx, ID)
		if err != nil {
			slog.Debug(fmt.Sprintf("could not retrieve revoked certificates: %v", err))
			return messages.ErrorMsg{
				Err: errors.Join(errors.New("could not parse CRL"), err),
			}
		}

		revokedCertificates := make([]x509.RevocationListEntry, len(certificates))
		for i, cert := range certificates {
			serialNumber := new(big.Int)
			serialNumber, ok := serialNumber.SetString(cert.SerialNumber, 10)
			if !ok {
				slog.Debug(fmt.Sprintf("could not parse serialNumber: %v", cert))
				return messages.ErrorMsg{
					Err: errors.Join(errors.New("could not parse serialNumber"), err),
				}
			}

			reasonCode := convertReasonCode(cert.RevocationReason)

			revokedCertificates[i] = x509.RevocationListEntry{
				SerialNumber:   serialNumber,
				RevocationTime: cert.RevocationDate,
				ReasonCode:     reasonCode,
			}
		}

		return messages.CRLResponseMsg{
			RevocationList: &x509.RevocationList{
				Issuer:                    pkix.Name{CommonName: aCN},
				ThisUpdate:                thisUpdate,
				NextUpdate:                nextUpdate,
				RevokedCertificateEntries: revokedCertificates,
			},
		}
	}
}

func convertReasonCode(reason crl.RevocationReason) int {
	switch reason {
	case crl.RevocationReasonUnspecified:
		return 0
	case crl.RevocationReasonKeyCompromise:
		return 1
	case crl.RevocationReasonCACompromise:
		return 2
	case crl.RevocationReasonAffiliationChanged:
		return 3
	case crl.RevocationReasonSuperseded:
		return 4
	case crl.RevocationReasonCessationOfOperation:
		return 5
	case crl.RevocationReasonCertificateHold:
		return 6
	case crl.RevocationReasonRemoveFromCRL:
		return 8
	case crl.RevocationReasonPriviledgeWithdrawn:
		return 9
	case crl.RevocationReasonAACompromise:
		return 10
	default:
		return 0
	}
}
