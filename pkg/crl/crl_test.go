package crl_test

import (
	"testing"

	"github.com/pimg/certguard/internal/adapter"
	"github.com/pimg/certguard/pkg/crl"
	"github.com/stretchr/testify/assert"
)

func TestCRL(t *testing.T) {
	adapter.GlobalCache = &adapter.MockCache{}

	testURL := "http://crl.quovadisglobal.com/pkioprivservg1.crl"

	res, err := crl.FetchRevocationList(testURL)
	assert.NoError(t, err)
	assert.NotNil(t, res)
}
