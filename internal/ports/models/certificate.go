package models

import (
	"crypto/x509"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pimg/certguard/internal/ports/models/styles"
)

type certificateKeyMap struct {
	Back key.Binding
	Quit key.Binding
	Home key.Binding
}

func (k *certificateKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Back, k.Quit, k.Home}
}

func (k *certificateKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Back, k.Quit, k.Home},
	}
}

var certificateKeys = certificateKeyMap{
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back to previous view"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Home: key.NewBinding(
		key.WithKeys("h"),
		key.WithHelp("h", "back to the main view"),
	),
}

type CertificateModel struct {
	keys        certificateKeyMap
	styles      *styles.Styles
	certificate x509.Certificate // TODO create custom struct and move parsing to domain model to do advanced parsing
}

func NewCertificateModel(cert *x509.Certificate) *CertificateModel {
	return &CertificateModel{
		keys:        certificateKeys,
		styles:      styles.DefaultStyles(),
		certificate: *cert,
	}
}

func (c *CertificateModel) Init() tea.Cmd {
	return nil
}

func (c *CertificateModel) Update(_ tea.Msg) (tea.Model, tea.Cmd) {
	return c, nil
}

func (c *CertificateModel) View() string {
	var s strings.Builder

	s.WriteString(c.styles.RevokedCertificateText.Render("Serialnumber: ") + c.certificate.SerialNumber.String())
	s.WriteString(c.styles.RevokedCertificateText.Render("CommonName: ") + c.certificate.Subject.CommonName)
	s.WriteString(c.parseCountry(c.certificate.Subject.Country))
	s.WriteString(c.styles.RevokedCertificateText.Render("Issuer: ") + c.certificate.Issuer.String())
	s.WriteString(c.styles.RevokedCertificateText.Render("NotBefore: ") + c.certificate.NotBefore.String())
	s.WriteString(c.styles.RevokedCertificateText.Render("NotAfter: ") + c.certificate.NotAfter.String())

	certInfo := c.styles.Text.Render(s.String())

	return lipgloss.JoinVertical(lipgloss.Top, certInfo)
}

func (c *CertificateModel) parseCountry(countries []string) string {
	if len(countries) == 0 {
		return ""
	}

	var s strings.Builder
	s.WriteString(c.styles.RevokedCertificateText.Render("Country: "))

	for _, country := range countries {
		s.WriteString(country)
	}

	return s.String()
}
