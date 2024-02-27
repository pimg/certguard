package commands

import (
	"errors"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pimg/certguard/internal/ports/models/messages"
	"github.com/pimg/certguard/pkg/crl"
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

func LoadCRL(path string) tea.Cmd {
	return func() tea.Msg {
		revocationList, err := crl.LoadRevocationList(strings.TrimSpace(path))
		if err != nil {
			return messages.ErrorMsg{
				Err: errors.Join(fmt.Errorf("could not load CRL from cache location: %s", path), err),
			}
		}
		return messages.CRLResponseMsg{
			RevocationList: revocationList,
		}
	}
}
