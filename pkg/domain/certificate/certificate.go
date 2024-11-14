package certificate

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
)

func ParsePEMCertificate(data []byte) ([]*x509.Certificate, error) {
	certificateChain := make([]*x509.Certificate, 0)
	for block, rest := pem.Decode(data); block != nil; block, rest = pem.Decode(rest) {
		switch block.Type {
		case "CERTIFICATE":
			cert, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				return nil, fmt.Errorf("failed to parse the PEM certificate: %v", err)
			}
			certificateChain = append(certificateChain, cert)
		default:
			return nil, errors.New("unsupported block type, only CERTIFICATE is supported")
		}
	}

	if len(certificateChain) == 0 {
		return nil, errors.New("failed to parse the PEM certificate")
	}

	return certificateChain, nil
}
