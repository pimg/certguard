package uri

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestURIValidator(t *testing.T) {
	t.Parallel()
	type testCase struct {
		input   string
		wantErr bool
	}

	testCases := []testCase{
		{
			input: "http://example.com",
		},
		{
			input: "https://example.com",
		},
		{
			input: "file:///tmp",
		},
		{
			input:   "invalid",
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.input, func(t *testing.T) {
			t.Parallel()
			err := ValidateURI(tc.input)

			if tc.wantErr {
				assert.NotNil(t, err)
				assert.ErrorContains(t, err, fmt.Sprintf("URI must start with either: 'http://', 'https://' or 'file://' the provided string: %s is not a valid URI:", tc.input))
			} else {
				assert.Nil(t, err)
			}
		})

	}
}
