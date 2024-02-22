package uri

import (
	"fmt"
	"net/url"
)

func ValidateURI(rawURI string) error {
	_, err := url.ParseRequestURI(rawURI)
	if err != nil {
		return fmt.Errorf("URI must start with either: 'http://', 'https://' or 'file://' the provided string: %s is not a valid URI: %s", rawURI, err)
	}

	return nil
}
