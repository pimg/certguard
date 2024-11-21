package models

import (
	"crypto/x509"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/tree"
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
	OSCP   key.Binding
}

func (k *certificateKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Search, k.OSCP, k.Back, k.Home}
}

func (k *certificateKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Search, k.Home},
		{k.OSCP},
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
	OSCP: key.NewBinding(
		key.WithKeys("o"),
		key.WithHelp("o", "perform OCSP request"),
	),
}

type CertificateModel struct {
	keys                 certificateKeyMap
	styles               *styles.Styles
	certificate          *x509.Certificate // TODO create custom struct and move parsing to domain model to do advanced parsing, also consider removing certificate as this can be obtained from the chiain
	certificateChain     []*x509.Certificate
	revocationInfo       *crl.RevokedCertificate
	foundOnCRL           *bool
	errorMsg             string
	OCSPStatus           string
	OCSPRevocationDate   time.Time
	OCSPRevocationReason string
	commands             *commands.Commands
}

func NewCertificateModel(cert *x509.Certificate, certificateChain []*x509.Certificate, cmds *commands.Commands) *CertificateModel {
	return &CertificateModel{
		keys:             certificateKeys,
		styles:           styles.DefaultStyles(),
		certificate:      cert,
		certificateChain: certificateChain,
		commands:         cmds,
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
		case "o":
			if len(c.certificate.OCSPServer) == 0 {
				c.errorMsg = "Certificate does not contain a OCSP Server"
				return c, cmd
			}

			if len(c.certificateChain) <= 1 {
				c.errorMsg = "Certificate does not contain a certificate chain, Issuer certificate missing"
				return c, cmd
			}
			cmd = c.commands.OCSPRequest(c.certificate, c.certificateChain[1], c.certificate.OCSPServer[0])
		}
	case messages.GetRevokedCertificateMsg:
		c.revocationInfo = msg.RevokedCertificate
		c.foundOnCRL = &msg.Found
	case messages.OCSPResponseMsg:
		c.OCSPStatus = msg.Status
		c.OCSPRevocationDate = msg.RevocationDate
		c.OCSPRevocationReason = msg.RevocationReason
	}
	return c, cmd
}

func (c *CertificateModel) View() string {
	var s strings.Builder
	s.WriteString(c.renderCertificateChain())

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

	if c.OCSPStatus != "" {
		s.WriteString("\n\nOCSP Response: \n")
		s.WriteString(c.styles.RevokedCertificateText.Render("OCSP Status: ") + c.OCSPStatus)
		if c.OCSPRevocationDate != (time.Time{}) {
			s.WriteString(c.styles.RevokedCertificateText.Render("Revocation Reason: ") + c.OCSPRevocationReason)
			s.WriteString(c.styles.RevokedCertificateText.Render("Revocation Date: ") + c.OCSPRevocationDate.Format(time.RFC3339))
		}
	}

	certInfo := c.styles.CertificateChain.Render(s.String())

	return lipgloss.JoinVertical(lipgloss.Top, certInfo)
}

func (c *CertificateModel) renderCertificateChain() string {
	certificate := c.certificateChain[0]
	t := tree.Root(c.styles.CertificateTitle.Render(certificate.Subject.CommonName)).
		Child(c.styles.CertificateText.Render("CommonName: ") + certificate.Subject.CommonName).
		Child(c.styles.CertificateText.Render("Serialnumber: ") + certificate.SerialNumber.String()).
		Child(c.styles.CertificateText.Render("DN: ") + parseDN(c.styles, certificate.Subject.String())).
		Child(parseCountry(c.styles, certificate.Subject.Country)).
		Child(c.styles.CertificateText.Render("Issuer: ") + certificate.Issuer.String()).
		Child(c.styles.CertificateText.Render("NotBefore: ") + certificate.NotBefore.String()).
		Child(c.styles.CertificateText.Render("NotAfter: ") + certificate.NotAfter.String())
	buildCertificateTree(c.styles, t, c.certificateChain[1:])
	return fmt.Sprint(t)
}

func buildCertificateTree(s *styles.Styles, t *tree.Tree, certificateChain []*x509.Certificate) {
	if len(certificateChain) == 0 {
		return
	}
	certificate := certificateChain[0]
	branch := tree.Root(s.CertificateTitle.Render(certificate.Subject.CommonName)).
		Child(s.CertificateText.Render("CommonName: ") + certificate.Subject.CommonName).
		Child(s.CertificateText.Render("Serialnumber: ") + certificate.SerialNumber.String()).
		Child(s.CertificateText.Render("DN: ") + parseDN(s, certificate.Subject.String())).
		Child(parseCountry(s, certificate.Subject.Country)).
		Child(s.CertificateText.Render("Issuer: ") + certificate.Issuer.String()).
		Child(s.CertificateText.Render("NotBefore: ") + certificate.NotBefore.String()).
		Child(s.CertificateText.Render("NotAfter: ") + certificate.NotAfter.String())
	t.Child(branch)
	buildCertificateTree(s, branch, certificateChain[1:])
}

func parseDN(s *styles.Styles, dn string) string {
	var str strings.Builder
	for _, k := range strings.Split(dn, ",") {
		sep := strings.Index(k, "=")
		str.WriteString(s.RevokedCertificateText.Render("\t"+k[:sep]+": ") + k[sep+1:])
	}

	return str.String()
}

func parseCountry(s *styles.Styles, countries []string) string {
	if len(countries) == 0 {
		return ""
	}

	var str strings.Builder
	str.WriteString(s.CertificateText.Render("Country: "))

	for _, country := range countries {
		str.WriteString(country)
	}

	return str.String()
}
