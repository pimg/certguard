package crl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCRL(t *testing.T) {
	testURL := "http://crl.quovadisglobal.com/pkioprivservg1.crl"

	res, err := FetchRevocationList(testURL)
	assert.NoError(t, err)
	assert.NotNil(t, res)
}
