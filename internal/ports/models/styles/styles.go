package styles

import "github.com/charmbracelet/lipgloss"

type Styles struct {
	InputField lipgloss.Style
}

func DefaultStyles() *Styles {
	return &Styles{
		InputField: lipgloss.NewStyle().BorderForeground(lipgloss.Color("36")).BorderStyle(lipgloss.NormalBorder()).Padding(1).Width(80),
	}
}
