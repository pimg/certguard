package crl

import (
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

func FetchRevocationList(revocationListURL string) (*x509.RevocationList, error) {
	client := http.Client{Timeout: 5 * time.Second}
	response, err := client.Get(revocationListURL)
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve CRL from revocationListURL: %s", revocationListURL)
	}
	defer response.Body.Close()

	if !strings.HasPrefix(response.Status, "2") {
		return nil, fmt.Errorf("server responded with a non 2xx status code: %s", response.Status)
	}

	rawCRL, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, errors.Join(err, fmt.Errorf("cannot parse HTTP response from %q", revocationListURL))
	}

	revocationList, err := ParseRevocationList(rawCRL)
	if err != nil {
		return nil, errors.Join(err, fmt.Errorf("cannot parse CRL from %q", revocationListURL))
	}
	return revocationList, nil
}

func ParseRevocationList(rawCRL []byte) (*x509.RevocationList, error) {
	revocationList, err := x509.ParseRevocationList(rawCRL)
	if err != nil {
		return nil, errors.Join(err, errors.New("cannot parse CRL from"))
	}

	return revocationList, nil
}
