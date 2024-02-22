package commands

import (
	"errors"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pimg/crl-inspector/internal/ports/models/messages"
	"github.com/pimg/crl-inspector/pkg/crl"
)

func GetCRL(requestURL string) tea.Cmd {
	return func() tea.Msg {
		revocationList, err := crl.FetchRevocationList(strings.TrimSpace(requestURL))
		if err != nil {
			return messages.ErrorMsg{
				Err: errors.Join(fmt.Errorf("could not download CRL with provided URL: %s", requestURL), err),
			}
		}
		return messages.CRLResponseMsg{
			RevocationList: revocationList,
		}
	}
}
