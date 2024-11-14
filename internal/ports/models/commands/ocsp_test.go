package commands

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/pimg/certguard/internal/ports/models/messages"
	"github.com/pimg/certguard/pkg/domain/crl"
	"github.com/stretchr/testify/assert"
)

func TestInvalidOCSPRequestCertificateIsNil(t *testing.T) {
	certRaw, err := os.ReadFile(filepath.Join("..", "..", "..", "..", "testing", "pki", "github.com-chain.pem"))
	assert.NoError(t, err)

	storage, err := crl.NewMockStorage()
	assert.NoError(t, err)

	cmds := NewCommands(storage)

	parseCertCmd := cmds.ParsePemCertficate(string(certRaw))
	assert.NotNil(t, parseCertCmd)

	certificateChainMsg := parseCertCmd()

	certificateChain := certificateChainMsg.(messages.PemCertificateMsg).CertificateChain
	issueCert := certificateChain[1]
	ocspServerURL := "http://example.com"

	cmd := cmds.OCSPRequest(nil, issueCert, ocspServerURL)

	ocspMsg := cmd()

	errMsg := ocspMsg.(messages.ErrorMsg)

	assert.Equal(t, "certificate is nil", errMsg.Err.Error())
}

func TestInvalidOCSPRequestCertificateIssuerIsNil(t *testing.T) {
	certRaw, err := os.ReadFile(filepath.Join("..", "..", "..", "..", "testing", "pki", "github.com-chain.pem"))
	assert.NoError(t, err)

	storage, err := crl.NewMockStorage()
	assert.NoError(t, err)

	cmds := NewCommands(storage)

	parseCertCmd := cmds.ParsePemCertficate(string(certRaw))
	assert.NotNil(t, parseCertCmd)

	certificateChainMsg := parseCertCmd()

	certificateChain := certificateChainMsg.(messages.PemCertificateMsg).CertificateChain
	cert := certificateChain[0]
	ocspServerURL := "http://example.com"

	cmd := cmds.OCSPRequest(cert, nil, ocspServerURL)

	ocspMsg := cmd()

	errMsg := ocspMsg.(messages.ErrorMsg)

	assert.Equal(t, "certificate Issuer is nil", errMsg.Err.Error())
}

func TestInvalidOCSPRequestInvalidURL(t *testing.T) {
	certRaw, err := os.ReadFile(filepath.Join("..", "..", "..", "..", "testing", "pki", "github.com-chain.pem"))
	assert.NoError(t, err)

	storage, err := crl.NewMockStorage()
	assert.NoError(t, err)

	cmds := NewCommands(storage)

	parseCertCmd := cmds.ParsePemCertficate(string(certRaw))
	assert.NotNil(t, parseCertCmd)

	certificateChainMsg := parseCertCmd()

	certificateChain := certificateChainMsg.(messages.PemCertificateMsg).CertificateChain
	cert := certificateChain[0]
	certChain := certificateChain[1]
	ocspServerURL := "invalidURL"

	cmd := cmds.OCSPRequest(cert, certChain, ocspServerURL)

	ocspMsg := cmd()

	errMsg := ocspMsg.(messages.ErrorMsg)

	assert.Equal(t, "could not validate OCSP server URL\nURI must start with either: 'http://', 'https://' or 'file://' the provided string: invalidURL is not a valid URI: parse \"invalidURL\": invalid URI for request", errMsg.Err.Error())
}

func TestOCSPInvalidResponse(t *testing.T) {
	certRaw, err := os.ReadFile(filepath.Join("..", "..", "..", "..", "testing", "pki", "github.com-chain.pem"))
	assert.NoError(t, err)

	storage, err := crl.NewMockStorage()
	assert.NoError(t, err)

	cmds := NewCommands(storage)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Invalid response")
	}))
	defer ts.Close()

	parseCertCmd := cmds.ParsePemCertficate(string(certRaw))
	assert.NotNil(t, parseCertCmd)

	certificateChainMsg := parseCertCmd()

	certificateChain := certificateChainMsg.(messages.PemCertificateMsg).CertificateChain
	cert := certificateChain[0]
	certChain := certificateChain[1]
	ocspServerURL := ts.URL

	cmd := cmds.OCSPRequest(cert, certChain, ocspServerURL)

	ocspMsg := cmd()

	errMsg := ocspMsg.(messages.ErrorMsg)

	assert.ErrorContains(t, errMsg.Err, "could not parse OCSP response for certificate")
}
