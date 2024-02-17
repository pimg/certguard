package styles

import "github.com/charmbracelet/lipgloss"

type Styles struct {
	InputField lipgloss.Style
	Title      lipgloss.Style
	Background lipgloss.Style
}

func DefaultStyles() *Styles {
	return &Styles{
		InputField: lipgloss.NewStyle().BorderForeground(lipgloss.Color("#458588")).BorderStyle(lipgloss.NormalBorder()).Padding(1).Width(80),
		Title: lipgloss.NewStyle().Bold(true).
			Foreground(lipgloss.Color("#D79921")).
			PaddingTop(2).
			PaddingBottom(2).
			PaddingLeft(2),
		Background: lipgloss.NewStyle().Background(lipgloss.Color("#282828")),
	}
}
