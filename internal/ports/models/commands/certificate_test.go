package commands

import (
	"math/big"
	"os"
	"path/filepath"
	"testing"

	"github.com/pimg/certguard/internal/ports/models/messages"
	"github.com/pimg/certguard/pkg/domain/crl"
	"github.com/stretchr/testify/assert"
)

func TestParseCertificate(t *testing.T) {
	certRaw, err := os.ReadFile(filepath.Join("..", "..", "..", "..", "testing", "pki", "org-on-crl.pem"))
	assert.NoError(t, err)

	storage, err := crl.NewMockStorage()
	assert.NoError(t, err)

	cmds := NewCommands(storage)

	cmd := cmds.ParsePemCertficate(string(certRaw))
	assert.NotNil(t, cmd)

	certMsg := cmd()
	assert.NotNil(t, certMsg)

	msg := certMsg.(messages.PemCertificateMsg)

	serialNum, success := new(big.Int).SetString("277698924469047062536476011533217874011933401810", 0)
	assert.True(t, success)
	assert.Equal(t, serialNum, msg.Certificate.SerialNumber)
}

func TestParseCertificateInvalid(t *testing.T) {
	storage, err := crl.NewMockStorage()
	assert.NoError(t, err)

	cmds := NewCommands(storage)

	cmd := cmds.ParsePemCertficate("this is not a valid certificate")
	assert.NotNil(t, cmd)

	msg := cmd()

	errMsg := msg.(messages.ErrorMsg)

	assert.Equal(t, "failed to parse certificate", errMsg.Err.Error())
}
