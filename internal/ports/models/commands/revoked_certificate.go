package commands

import (
	"context"
	"crypto/x509"
	"crypto/x509/pkix"
	"errors"
	"log"
	"math/big"
	"net/url"
	"strconv"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pimg/certguard/internal/ports/models/messages"
	"github.com/pimg/certguard/pkg/domain/crl"
)

type GetRevokedCertificatesArgs struct {
	ID         string
	CN         string
	ThisUpdate string
	NextUpdate string
	URL        string
}

func (c *Commands) GetRevokedCertificates(args *GetRevokedCertificatesArgs) tea.Cmd {
	log.Printf("getting revoked certificate from storage with ID: %s", args.ID)
	ctx := context.Background()
	return func() tea.Msg {
		ID, err := strconv.ParseInt(args.ID, 10, 64)
		if err != nil {
			log.Printf("%s is not a valid ID, %v", args.ID, err)
			return messages.ErrorMsg{
				Err: errors.Join(errors.New("could not parse ID"), err),
			}
		}

		thisUpdate, err := time.Parse(time.DateOnly, args.ThisUpdate)
		if err != nil {
			log.Printf("%s is not a valid time, %v", args.ThisUpdate, err)
			return messages.ErrorMsg{
				Err: errors.Join(errors.New("could not parse time of ThisUpdate"), err),
			}
		}

		nextUpdate, err := time.Parse(time.DateOnly, args.NextUpdate)
		if err != nil {
			log.Printf("%s is not a valid time, %v", args.NextUpdate, err)
			return messages.ErrorMsg{
				Err: errors.Join(errors.New("could not parse time of NextUpdate"), err),
			}
		}

		var URL *url.URL
		if args.URL != "" {
			URL, err = url.ParseRequestURI(args.URL)
			if err != nil {
				log.Printf("%s is not a valid URL, %v", args.URL, err)
				return messages.ErrorMsg{
					Err: errors.Join(errors.New("could not parse URL"), err),
				}
			}
		}

		certificates, err := c.storage.Repository.FindRevokedCertificates(ctx, ID)
		if err != nil {
			log.Printf("could not retrieve revoked certificates: %v", err)
			return messages.ErrorMsg{
				Err: errors.Join(errors.New("could not parse CRL"), err),
			}
		}

		revokedCertificates := make([]x509.RevocationListEntry, len(certificates))
		for i, cert := range certificates {
			serialNumber := new(big.Int)
			serialNumber, ok := serialNumber.SetString(cert.SerialNumber, 10)
			if !ok {
				log.Printf("could not parse serialNumber: %v", cert)
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
				Issuer:                    pkix.Name{CommonName: args.CN},
				ThisUpdate:                thisUpdate,
				NextUpdate:                nextUpdate,
				RevokedCertificateEntries: revokedCertificates,
			},
			URL: URL,
		}
	}
}

func (c *Commands) Search(serialnumber string) tea.Cmd {
	log.Printf("search stored CRLs for serialnumber: %s", serialnumber)
	ctx := context.Background()
	return func() tea.Msg {
		revokedCertificate, err := c.storage.Repository.FindRevokedCertificate(ctx, serialnumber)
		if err != nil {
			log.Printf("could not perform find action on serialnumber: %s", serialnumber)
			return messages.ErrorMsg{
				Err: errors.Join(errors.New("could not perform find action on serial number"), err),
			}
		}

		if revokedCertificate == nil {
			log.Printf("no revoked certificate for serialnumber: %s", serialnumber)
			return messages.GetRevokedCertificateMsg{
				Found: false,
			}
		}

		return messages.GetRevokedCertificateMsg{
			RevokedCertificate: revokedCertificate,
			Found:              true,
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
