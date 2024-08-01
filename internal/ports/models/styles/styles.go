package styles

import (
	"github.com/charmbracelet/lipgloss"
)

type Styles struct {
	InputField             lipgloss.Style
	Title                  lipgloss.Style
	Background             lipgloss.Style
	ErrorMessages          lipgloss.Style
	Text                   lipgloss.Style
	BaseText               lipgloss.Style
	RevokedCertificateText lipgloss.Style
	CRLText                lipgloss.Style
	FilePickerFile         lipgloss.Style
	FilePickerCurrent      lipgloss.Style
	ListComponentTitle     lipgloss.Color
	WarningText            lipgloss.Style
}

func DefaultStyles() *Styles {
	return &Styles{
		InputField: lipgloss.NewStyle().BorderForeground(lipgloss.Color("#83A598")).BorderStyle(lipgloss.NormalBorder()).Padding(1).Width(80),
		Title: lipgloss.NewStyle().Bold(true).
			Foreground(lipgloss.Color("#EBDBB2")).
			Background(lipgloss.Color("#83A598")).
			Width(80).
			PaddingTop(1).
			PaddingBottom(1).
			PaddingLeft(1),
		Background:             lipgloss.NewStyle().Background(lipgloss.Color("#282828")),
		ErrorMessages:          lipgloss.NewStyle().Background(lipgloss.Color("#FB4934")).BorderForeground(lipgloss.Color("#FB4934")).BorderStyle(lipgloss.NormalBorder()).Width(80).Padding(1),
		Text:                   lipgloss.NewStyle().Foreground(lipgloss.Color("#B8BB26")).Padding(1).Width(80),
		RevokedCertificateText: lipgloss.NewStyle().Foreground(lipgloss.Color("#B8BB26")).PaddingTop(1).PaddingLeft(1).Width(20),
		CRLText:                lipgloss.NewStyle().Foreground(lipgloss.Color("#B8BB26")).PaddingTop(1).PaddingLeft(1).Width(25),
		BaseText:               lipgloss.NewStyle().Foreground(lipgloss.Color("#B8BB26")).PaddingLeft(1).Width(48),
		FilePickerFile:         lipgloss.NewStyle().Foreground(lipgloss.Color("#83A598")),
		FilePickerCurrent:      lipgloss.NewStyle().Foreground(lipgloss.Color("#B8BB26")),
		ListComponentTitle:     "#83A598",
		WarningText:            lipgloss.NewStyle().Foreground(lipgloss.Color("#FABD2F")),
	}
}
