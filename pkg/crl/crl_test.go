package crl_test

import (
	"testing"

	"github.com/pimg/certguard/internal/adapter"
	"github.com/pimg/certguard/pkg/crl"
	"github.com/stretchr/testify/assert"
)

type mockCache struct {
	filename    string
	fileContent []byte
}

func (m *mockCache) Write(filename string, fileContent []byte) error {
	m.filename = filename
	m.fileContent = fileContent

	return nil
}

func TestCRL(t *testing.T) {
	adapter.GlobalCache = &mockCache{}

	testURL := "http://crl.quovadisglobal.com/pkioprivservg1.crl"

	res, err := crl.FetchRevocationList(testURL)
	assert.NoError(t, err)
	assert.NotNil(t, res)
}
