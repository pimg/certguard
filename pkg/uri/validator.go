package uri

import (
	"fmt"
	"net/url"
)

func ValidateURI(rawURI string) (*url.URL, error) {
	url, err := url.ParseRequestURI(rawURI)
	if err != nil {
		return nil, fmt.Errorf("URI must start with either: 'http://', 'https://' or 'file://' the provided string: %s is not a valid URI: %s", rawURI, err)
	}

	return url, nil
}
