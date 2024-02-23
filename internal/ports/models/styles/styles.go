package styles

import "github.com/charmbracelet/lipgloss"

type Styles struct {
	InputField    lipgloss.Style
	Title         lipgloss.Style
	Background    lipgloss.Style
	ErrorMessages lipgloss.Style
	Text          lipgloss.Style
}

func DefaultStyles() *Styles {
	return &Styles{
		InputField: lipgloss.NewStyle().BorderForeground(lipgloss.Color("#83A598")).BorderStyle(lipgloss.NormalBorder()).Padding(1).Width(80),
		Title: lipgloss.NewStyle().Bold(true).
			Foreground(lipgloss.Color("#FABD2F")).
			PaddingTop(2).
			PaddingBottom(2).
			PaddingLeft(2),
		Background:    lipgloss.NewStyle().Background(lipgloss.Color("#282828")),
		ErrorMessages: lipgloss.NewStyle().Background(lipgloss.Color("#FB4934")).BorderForeground(lipgloss.Color("#FB4934")).BorderStyle(lipgloss.NormalBorder()).Width(80).Padding(1),
		Text:          lipgloss.NewStyle().Foreground(lipgloss.Color("#B8BB26")).Padding(1).Width(80),
	}
}
