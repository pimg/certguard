package models

import (
	"crypto/x509"
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pimg/certguard/internal/ports/models/styles"
)

// keyMap defines a set of keybindings. To work for help it must satisfy
// key.Map. It could also very easily be a map[string]key.Binding.
type listKeyMap struct {
	Back key.Binding
	Quit key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k *listKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Back, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k *listKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Back, k.Quit},
	}
}

var listKeys = listKeyMap{
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back to main view"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

type ListModel struct {
	keys   listKeyMap
	help   help.Model
	styles *styles.Styles
	crl    *x509.RevocationList
}

func NewListModel(crl *x509.RevocationList) ListModel {
	return ListModel{
		keys:   listKeys,
		help:   help.New(),
		styles: styles.DefaultStyles(),
		crl:    crl,
	}
}

func (l ListModel) Init() tea.Cmd {
	return nil
}

func (l ListModel) Update(_ tea.Msg) (tea.Model, tea.Cmd) {
	return l, nil
}

func (l ListModel) View() string {
	issuer := fmt.Sprintf("CRL Issuer          : %s", l.crl.Issuer)
	updatedAt := fmt.Sprintf("Updated At          : %s", l.crl.ThisUpdate)
	nextUpdate := fmt.Sprintf("Next Update         : %s", l.crl.NextUpdate)
	revokedCertCount := fmt.Sprintf("Revoked Certificates: %d", len(l.crl.RevokedCertificateEntries))

	crlInfo := l.styles.Text.Render(
		fmt.Sprintf("%s\n%s\n%s\n%s", issuer, updatedAt, nextUpdate, revokedCertCount),
	)
	return crlInfo
}
