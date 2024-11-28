package models

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pimg/certguard/internal/ports/models/styles"
	"github.com/pimg/certguard/pkg/domain/crl"
)

type revokedCertKeyMap struct {
	Back key.Binding
	Quit key.Binding
}

func (k *revokedCertKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Back, k.Quit}
}

func (k *revokedCertKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Back, k.Quit},
	}
}

var revokedCertificateKeys = revokedCertKeyMap{
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back to previous view"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

type RevokedCertificateModel struct {
	serialnumber, revocationReason, revocationDate string
	styles                                         *styles.Styles
	keys                                           revokedCertKeyMap
}

func NewRevokedCertificateModel(serialnumber, revocationReason, revocationDate string) *RevokedCertificateModel {
	return &RevokedCertificateModel{
		serialnumber:     serialnumber,
		revocationReason: revocationReason,
		revocationDate:   revocationDate,
		keys:             revokedCertificateKeys,
		styles:           styles.Theme,
	}
}

func (r *RevokedCertificateModel) Init() tea.Cmd {
	return nil
}

func (r *RevokedCertificateModel) Update(_ tea.Msg) (tea.Model, tea.Cmd) {
	return r, nil
}

func (r *RevokedCertificateModel) View() string {
	serialnumber := r.styles.RevokedCertificateText.Render("Serialnumber: ") + r.serialnumber
	revocationDate := r.styles.RevokedCertificateText.Render("Revocation date: ") + r.revocationDate

	var revocationReason string
	revocationReasonCode, err := strconv.Atoi(r.revocationReason)
	if err != nil {
		revocationReason = "Unknown revocation reason"
	} else {
		revocationReason = r.styles.RevokedCertificateText.Render("Revocation reason: ") + crl.RevocationReasons[revocationReasonCode].String()
	}

	return fmt.Sprintf("%s%s%s", serialnumber, revocationDate, revocationReason)
}
