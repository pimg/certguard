package certificate

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
)

func ParsePEMCertificate(data []byte) (*x509.Certificate, error) {
	block, restBlocks := pem.Decode(data)
	if len(restBlocks) > 0 {
		return nil, errors.New("the PEM files contains multiple certificates, only one is supported")
	}

	if block.Type != "CERTIFICATE" {
		return nil, errors.New("the PEM files must contain a certificate")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the PEM certificate: %v", err)
	}

	return cert, nil
}
