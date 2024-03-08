package commands

import (
	"errors"
	"fmt"
	"log"
	"log/slog"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pimg/certguard/internal/adapter"
	"github.com/pimg/certguard/internal/ports/models/messages"
	"github.com/pimg/certguard/pkg/crl"
)

func GetCRL(requestURL string) tea.Cmd {
	slog.Debug("requesting CRL from: " + requestURL)
	return func() tea.Msg {
		revocationListURL := strings.TrimSpace(requestURL)
		revocationList, err := crl.FetchRevocationList(revocationListURL)
		if err != nil {
			errorMsg := fmt.Errorf("could not download CRL with provided URL: %s", requestURL)
			log.Print(errorMsg.Error())
			return messages.ErrorMsg{
				Err: errors.Join(errorMsg, err),
			}
		}

		filename := revocationListURL[strings.LastIndex(revocationListURL, "/"):]

		err = adapter.GlobalCache.Write(filename, revocationList.Raw)
		if err != nil {
			errorMsg := fmt.Errorf("cannot write the CRL to the cache: %s", filename)
			log.Print(errorMsg.Error())
			return messages.ErrorMsg{Err: errors.Join(err, errorMsg)}
		}

		return messages.CRLResponseMsg{
			RevocationList: revocationList,
		}
	}
}

func LoadCRL(path string) tea.Cmd {
	return func() tea.Msg {
		rawCRL, err := adapter.GlobalCache.Read(path)
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
		return messages.CRLResponseMsg{
			RevocationList: revocationList,
		}
	}
}
