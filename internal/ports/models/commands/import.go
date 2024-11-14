package commands

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pimg/certguard/internal/ports/models/messages"
	"github.com/pimg/certguard/pkg/crl"
	"github.com/pimg/certguard/pkg/domain/certificate"
	domain_crl "github.com/pimg/certguard/pkg/domain/crl"
)

func (c *Commands) ImportFile(path string) tea.Cmd {
	log.Printf("loading CRL from path: %s", path)
	ctx := context.Background()
	return func() tea.Msg {
		rawFile, err := os.ReadFile(path)
		if err != nil {
			errorMsg := fmt.Errorf("could not load CRL from import location: %s", path)
			log.Println(errorMsg.Error())
			return messages.ErrorMsg{
				Err: errors.Join(errorMsg, err),
			}
		}

		switch filepath.Ext(path) {
		case ".crl":
			log.Println("importing CRL based on file extension")
			revocationList, err := crl.ParseRevocationList(rawFile)
			if err != nil {
				log.Println("could not parse CRL")
				return messages.ErrorMsg{
					Err: errors.Join(errors.New("could not parse CRL"), err),
				}
			}

			err = domain_crl.Process(ctx, nil, revocationList, c.storage)
			if err != nil {
				return nil
			}

			return messages.CRLResponseMsg{
				RevocationList: revocationList,
			}
		default:
			log.Println("importing Certificate based on file extension")
			certificateChain, err := certificate.ParsePEMCertificate(rawFile)
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
}
