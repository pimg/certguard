package commands

import (
	"errors"
	"log"
	"slices"

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

		slices.Reverse(certificateChain)
		log.Println("reversed certificate chain")
		return messages.PemCertificateMsg{
			Certificate:      certificateChain[len(certificateChain)-1],
			CertificateChain: certificateChain,
		}
	}
}
