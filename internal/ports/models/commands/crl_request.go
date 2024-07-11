package commands

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/url"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pimg/certguard/internal/ports/models/messages"
	"github.com/pimg/certguard/pkg/crl"
	domain_crl "github.com/pimg/certguard/pkg/domain/crl"
)

func GetCRL(url *url.URL) tea.Cmd {
	slog.Debug("requesting CRL from: " + url.String())
	ctx := context.Background()
	return func() tea.Msg {
		revocationListURL := strings.TrimSpace(url.String())
		revocationList, err := crl.FetchRevocationList(revocationListURL)
		if err != nil {
			errorMsg := fmt.Errorf("could not download CRL with provided URL: %s", url.String())
			log.Print(errorMsg.Error())
			return messages.ErrorMsg{
				Err: errors.Join(errorMsg, err),
			}
		}

		err = domain_crl.Process(ctx, url, revocationList, domain_crl.GlobalStorage)
		if err != nil {
			return nil
		}

		return messages.CRLResponseMsg{
			RevocationList: revocationList,
		}
	}
}

func LoadCRL(path string) tea.Cmd {
	ctx := context.Background()
	return func() tea.Msg {
		rawCRL, err := os.ReadFile(path)
		if err != nil {
			errorMsg := fmt.Errorf("could not load CRL from cache location: %s", path)
			log.Print(errorMsg.Error())
			return messages.ErrorMsg{
				Err: errors.Join(errorMsg, err),
			}
		}

		revocationList, err := crl.ParseRevocationList(rawCRL)
		if err != nil {
			log.Print("could not parse CRL")
			return messages.ErrorMsg{
				Err: errors.Join(errors.New("could not parse CRL"), err),
			}
		}

		err = domain_crl.Process(ctx, nil, revocationList, domain_crl.GlobalStorage)
		if err != nil {
			return nil
		}

		return messages.CRLResponseMsg{
			RevocationList: revocationList,
		}
	}
}

func GetCRLsFromStore() tea.Msg {
	ctx := context.Background()
	cRLs, err := domain_crl.GlobalStorage.Repository.List(ctx)
	if err != nil {
		return nil
	}

	return messages.ListCRLsResponseMsg{
		CRLs: cRLs,
	}
}
