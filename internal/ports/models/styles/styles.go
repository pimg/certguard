package styles

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/pimg/certguard/internal/ports/models/styles/colors"
)

var Theme *Styles

type Styles struct {
	InputField             lipgloss.Style
	TextArea               lipgloss.Style
	Title                  lipgloss.Style
	Background             lipgloss.Style
	ErrorMessages          lipgloss.Style
	Text                   lipgloss.Style
	BaseText               lipgloss.Style
	BaseMenuText           lipgloss.Style
	RevokedCertificateText lipgloss.Style
	CRLText                lipgloss.Style
	FilePickerFile         lipgloss.Style
	FilePickerCurrent      lipgloss.Style
	ListComponentTitle     lipgloss.Color
	WarningText            lipgloss.Style
	CertificateChain       lipgloss.Style
	CertificateTitle       lipgloss.Style
	CertificateText        lipgloss.Style
	CertificateWarning     lipgloss.Style
}

func gruvboxTheme() *colors.ThemeColors {
	return colors.SetThemeColors(&colors.ThemeColorArgs{
		MainBanner:    "#83A598",
		Background:    "#282828",
		HighlightText: "#B8BB26",
		Text:          "#EBDBB2",
		WarningText:   "#FABD2F",
		ErrorText:     "#FB4934",
	})
}

func draculaTheme() *colors.ThemeColors {
	return colors.SetThemeColors(&colors.ThemeColorArgs{
		MainBanner:    "#bd93f9",
		Background:    "#282a36",
		HighlightText: "#50fa7b",
		Text:          "#f8f8f2",
		WarningText:   "#f1fa8c",
		ErrorText:     "#ff5555",
	})
}

func NewStyles(themeName string) {
	var themeColors *colors.ThemeColors
	switch themeName {
	case "gruvbox":
		themeColors = gruvboxTheme()
	case "dracula":
		themeColors = draculaTheme()
	default:
		themeColors = gruvboxTheme()
	}
	Theme = &Styles{
		InputField: lipgloss.NewStyle().BorderForeground(themeColors.MainBanner).BorderStyle(lipgloss.NormalBorder()).Padding(1).Width(80),
		TextArea:   lipgloss.NewStyle().BorderForeground(themeColors.MainBanner).BorderStyle(lipgloss.NormalBorder()).Padding(1).Width(78),
		Title: lipgloss.NewStyle().Bold(true).
			Foreground(themeColors.Text).
			Background(themeColors.MainBanner).
			Width(900).
			PaddingTop(1).
			PaddingBottom(1).
			PaddingLeft(1),
		Background:             lipgloss.NewStyle().Background(themeColors.Background),
		ErrorMessages:          lipgloss.NewStyle().Background(themeColors.ErrorText).BorderForeground(themeColors.ErrorText).BorderStyle(lipgloss.NormalBorder()).Width(80).Padding(1),
		Text:                   lipgloss.NewStyle().Foreground(themeColors.HighlightText).Padding(1).Width(80),
		RevokedCertificateText: lipgloss.NewStyle().Foreground(themeColors.HighlightText).PaddingTop(1).PaddingLeft(1).Width(20),
		CRLText:                lipgloss.NewStyle().Foreground(themeColors.HighlightText).PaddingTop(1).PaddingLeft(1).Width(25),
		BaseText:               lipgloss.NewStyle().Foreground(themeColors.HighlightText).PaddingLeft(1).Width(48),
		BaseMenuText:           lipgloss.NewStyle().Foreground(themeColors.HighlightText).PaddingLeft(1).Width(60),
		FilePickerFile:         lipgloss.NewStyle().Foreground(themeColors.MainBanner),
		FilePickerCurrent:      lipgloss.NewStyle().Foreground(themeColors.HighlightText),
		ListComponentTitle:     themeColors.MainBanner,
		WarningText:            lipgloss.NewStyle().Foreground(themeColors.WarningText),
		CertificateWarning:     lipgloss.NewStyle().Foreground(themeColors.WarningText).PaddingTop(1).PaddingLeft(1).PaddingBottom(1),
		CertificateChain:       lipgloss.NewStyle().PaddingTop(1).PaddingLeft(1),
		CertificateTitle:       lipgloss.NewStyle().Foreground(themeColors.HighlightText),
		CertificateText:        lipgloss.NewStyle().PaddingLeft(1).Width(20).Foreground(themeColors.MainBanner),
	}
}
