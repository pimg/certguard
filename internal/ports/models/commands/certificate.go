package commands

import (
	"errors"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pimg/certguard/internal/ports/models/messages"
	"github.com/pimg/certguard/pkg/domain/certificate"
)

func (c *Commands) ParsePemCertficate(pem string) tea.Cmd {
	return func() tea.Msg {
		certificateChain, err := certificate.ParsePEMCertificate([]byte(pem))
		if err != nil {
			log.Printf("failed to parse certificate: %s", err)
			return messages.ErrorMsg{
				Err: errors.New("failed to parse certificate"),
			}
		}

		return messages.PemCertificateMsg{
			Certificate:      certificateChain[0],
			CertificateChain: certificateChain,
		}
	}
}
