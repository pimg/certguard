package commands

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/pimg/certguard/internal/ports/models/messages"
	"github.com/pimg/certguard/pkg/domain/crl"
	"github.com/stretchr/testify/assert"
)

func TestImportCRL(t *testing.T) {
	storage, err := crl.NewMockStorage()
	assert.NoError(t, err)

	cmds := NewCommands(storage)

	msg := cmds.ImportFile(filepath.Join("..", "..", "..", "..", "testing", "pki", "ca.crl"))()

	crlMsg := msg.(messages.CRLResponseMsg)

	assert.Len(t, crlMsg.RevocationList.RevokedCertificateEntries, 1)
}

func TestImportPEM(t *testing.T) {
	storage, err := crl.NewMockStorage()
	assert.NoError(t, err)

	cmds := NewCommands(storage)

	msg := cmds.ImportFile(filepath.Join("..", "..", "..", "..", "testing", "pki", "org-on-crl.pem"))()

	pemMsg := msg.(messages.PemCertificateMsg)

	assert.NotEmpty(t, pemMsg)
	assert.Equal(t, "277698924469047062536476011533217874011933401810", pemMsg.Certificate.SerialNumber.String())
}

func TestImportCRLInvalidPath(t *testing.T) {
	storage, err := crl.NewMockStorage()
	assert.NoError(t, err)

	cmds := NewCommands(storage)

	msg := cmds.ImportFile(filepath.Join("..", "..", "..", "..", "testing", "pki", "idonotexist.crl"))()

	errMsg := msg.(messages.ErrorMsg)

	assert.ErrorContains(t, errMsg.Err, "could not load CRL from import location")
}

func TestImportMalformedCertificate(t *testing.T) {
	storage, err := crl.NewMockStorage()
	assert.NoError(t, err)

	cmds := NewCommands(storage)

	msg := cmds.ImportFile(filepath.Join("..", "..", "..", "..", "testing", "pki", "malformed-certificate.pem"))()

	errMsg := msg.(messages.ErrorMsg)
	assert.ErrorContains(t, errMsg.Err, "failed to parse certificate")
}

func TestImportMalformedCRL(t *testing.T) {
	storage, err := crl.NewMockStorage()
	assert.NoError(t, err)

	cmds := NewCommands(storage)

	msg := cmds.ImportFile(filepath.Join("..", "..", "..", "..", "testing", "pki", "malformed.crl"))()

	errMsg := msg.(messages.ErrorMsg)
	assert.ErrorContains(t, errMsg.Err, "could not parse CRL\nx509: malformed crl\ncannot parse CRL from")
}

func TestGetRevokedCertificatesNoRevokedCertificatesFound(t *testing.T) {
	storage, err := crl.NewMockStorage()
	assert.NoError(t, err)

	cmds := NewCommands(storage)

	msg := cmds.GetRevokedCertificates(&GetRevokedCertificatesArgs{
		ID:         "1",
		CN:         "testCN",
		ThisUpdate: time.Now().Format(time.DateOnly),
		NextUpdate: time.Now().Format(time.DateOnly),
		URL:        "http://example.com",
	})()

	assert.NotNil(t, msg)
}
