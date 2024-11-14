package models

import (
	"crypto/x509"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pimg/certguard/internal/ports/models/commands"
	"github.com/pimg/certguard/internal/ports/models/messages"
	"github.com/pimg/certguard/internal/ports/models/styles"
	"github.com/pimg/certguard/pkg/domain/crl"
)

type certificateKeyMap struct {
	Back   key.Binding
	Quit   key.Binding
	Home   key.Binding
	Search key.Binding
}

func (k *certificateKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Search, k.Back, k.Quit, k.Home}
}

func (k *certificateKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Search, k.Home},
		{k.Back, k.Quit},
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
	Search: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "search CRLs with this certificate"),
	),
}

type CertificateModel struct {
	keys           certificateKeyMap
	styles         *styles.Styles
	certificate    *x509.Certificate // TODO create custom struct and move parsing to domain model to do advanced parsing
	revocationInfo *crl.RevokedCertificate
	foundOnCRL     *bool
	errorMsg       string
	commands       *commands.Commands
}

func NewCertificateModel(cert *x509.Certificate, cmds *commands.Commands) *CertificateModel {
	return &CertificateModel{
		keys:        certificateKeys,
		styles:      styles.DefaultStyles(),
		certificate: cert,
		commands:    cmds,
	}
}

func (c *CertificateModel) Init() tea.Cmd {
	return nil
}

func (c *CertificateModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case messages.ErrorMsg:
		c.errorMsg = msg.Err.Error()
	case tea.KeyMsg:
		switch msg.String() {
		case "s":
			cmd = c.commands.Search(c.certificate.SerialNumber.String())
			return c, cmd
		}
	case messages.GetRevokedCertificateMsg:
		c.revocationInfo = msg.RevokedCertificate
		c.foundOnCRL = &msg.Found
	}
	return c, cmd
}

func (c *CertificateModel) View() string {
	var s strings.Builder

	s.WriteString(c.styles.RevokedCertificateText.Render("Serialnumber: ") + c.certificate.SerialNumber.String())
	s.WriteString(c.styles.RevokedCertificateText.Render("CommonName: ") + c.certificate.Subject.CommonName)
	s.WriteString(c.styles.RevokedCertificateText.Render("DN: ") + c.parseDN(c.certificate.Subject.String()))
	s.WriteString(c.parseCountry(c.certificate.Subject.Country))
	s.WriteString(c.styles.RevokedCertificateText.Render("Issuer: ") + c.certificate.Issuer.String())
	s.WriteString(c.styles.RevokedCertificateText.Render("NotBefore: ") + c.certificate.NotBefore.String())
	s.WriteString(c.styles.RevokedCertificateText.Render("NotAfter: ") + c.certificate.NotAfter.String())

	if c.errorMsg != "" {
		s.WriteString("\n\n\n" + c.errorMsg)
	}

	if c.foundOnCRL != nil {
		s.WriteString("\n\nRevocation Info: \n")
		if *c.foundOnCRL && c.revocationInfo != nil {
			s.WriteString(c.styles.RevokedCertificateText.Render("Certificate is revoked!"))
			s.WriteString(c.styles.RevokedCertificateText.Render("Revocation Reason: ") + c.revocationInfo.RevocationReason.String())
			s.WriteString(c.styles.RevokedCertificateText.Render("Revocation Date: ") + c.revocationInfo.RevocationDate.String())
			s.WriteString(c.styles.RevokedCertificateText.Render("Revoked by: ") + c.revocationInfo.RevokedBy)
		}

		if !*c.foundOnCRL {
			s.WriteString(c.styles.Text.Render("Certificate is not found on any of the stored CRLs"))
		}
	}

	certInfo := c.styles.Text.Render(s.String())

	return lipgloss.JoinVertical(lipgloss.Top, certInfo)
}

func (c *CertificateModel) parseDN(dn string) string {
	var s strings.Builder
	for _, k := range strings.Split(dn, ",") {
		sep := strings.Index(k, "=")
		s.WriteString(c.styles.RevokedCertificateText.Render("\t"+k[:sep]+": ") + k[sep+1:])
	}

	return s.String()
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
