package commands

import (
	"bytes"
	"crypto"
	"crypto/x509"
	"errors"
	"io"
	"log"
	"net/http"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pimg/certguard/internal/ports/models/messages"
	"github.com/pimg/certguard/pkg/uri"
	"golang.org/x/crypto/ocsp"
)

func (c *Commands) OCSPRequest(cert, issuerCert *x509.Certificate, ocspServerURL string) tea.Cmd {
	opts := ocsp.RequestOptions{
		Hash: crypto.SHA256,
	}

	return func() tea.Msg {
		if cert == nil {
			log.Printf("certificate is nil")
			return messages.ErrorMsg{
				Err: errors.New("certificate is nil"),
			}
		}

		if issuerCert == nil {
			log.Printf("certificate Issuer is nil")
			return messages.ErrorMsg{
				Err: errors.New("certificate Issuer is nil"),
			}
		}

		OCSPServerURL, err := uri.ValidateURI(ocspServerURL)
		if err != nil {
			log.Printf("could not validate OCSP server URL: %s, err: %v", ocspServerURL, err)
			return messages.ErrorMsg{
				Err: errors.Join(errors.New("could not validate OCSP server URL"), err),
			}
		}

		log.Printf("Querying OCSP server URL: %s, for certificate: %s", ocspServerURL, cert.SerialNumber.String())

		buffer, err := ocsp.CreateRequest(cert, issuerCert, &opts)
		if err != nil {
			log.Printf("could not create OCSP request for certificate: %s", cert.SerialNumber.String())
			return messages.ErrorMsg{
				Err: errors.Join(errors.New("could not create OCSP request"), err),
			}
		}

		httpRequest, err := http.NewRequest(http.MethodPost, OCSPServerURL.String(), bytes.NewReader(buffer))
		if err != nil {
			log.Printf("could not create OCSP request for certificate: %s, err: %v", cert.SerialNumber.String(), err)
			return errors.Join(errors.New("could not create OCSP request"), err)
		}

		httpRequest.Header.Add("Content-Type", "application/ocsp-request")
		httpRequest.Header.Add("Accept", "application/ocsp-response")
		httpRequest.Header.Add("Host", OCSPServerURL.Hostname())

		httpClient := &http.Client{}
		httpResponse, err := httpClient.Do(httpRequest)
		if err != nil {
			log.Printf("could not send OCSP request for certificate: %s, err: %v", cert.SerialNumber.String(), err)
			return messages.ErrorMsg{
				Err: errors.Join(errors.New("could not send OCSP request"), err),
			}
		}
		defer httpResponse.Body.Close()
		OCSPResponseRaw, err := io.ReadAll(httpResponse.Body)
		if err != nil {
			log.Printf("could not read OCSP response for certificate: %s, err: %v", cert.SerialNumber.String(), err)
			return messages.ErrorMsg{
				Err: errors.Join(errors.New("could not read OCSP response for certificate"), err),
			}
		}

		OCSPResponse, err := ocsp.ParseResponseForCert(OCSPResponseRaw, cert, issuerCert)
		if err != nil {
			log.Printf("could not parse OCSP response for certificate: %s, err: %v", cert.SerialNumber, err)
			return messages.ErrorMsg{
				Err: errors.Join(errors.New("could not parse OCSP response for certificate"), err),
			}
		}

		switch OCSPResponse.Status {
		case ocsp.Good:
			return messages.OCSPResponseMsg{Status: "Good"}
		case ocsp.Revoked:
			revocationReason := parseRevocationReason(OCSPResponse.RevocationReason)
			return messages.OCSPResponseMsg{Status: "Revoked", RevocationDate: OCSPResponse.RevokedAt, RevocationReason: revocationReason}
		case ocsp.Unknown:
			return messages.OCSPResponseMsg{Status: "Unknown"}
		default:
			return messages.OCSPResponseMsg{Status: "Unknown"}
		}
	}
}

func parseRevocationReason(reason int) string {
	switch reason {
	case ocsp.Unspecified:
		return "unspecified"
	case ocsp.KeyCompromise:
		return "key compromise"
	case ocsp.CACompromise:
		return "CA compromise"
	case ocsp.AffiliationChanged:
		return "affiliation changed"
	case ocsp.Superseded:
		return "superseded"
	case ocsp.CessationOfOperation:
		return "cessation of operation"
	case ocsp.CertificateHold:
		return "certificate hold"
	case ocsp.RemoveFromCRL:
		return "remove from CRL"
	case ocsp.PrivilegeWithdrawn:
		return "privilege-withdrawn"
	case ocsp.AACompromise:
		return "AA compromise"
	default:
		return "unknown"
	}
}
