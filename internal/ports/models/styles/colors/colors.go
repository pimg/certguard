package colors

import "github.com/charmbracelet/lipgloss"

type ThemeColors struct {
	MainBanner    lipgloss.Color
	Background    lipgloss.Color
	HighlightText lipgloss.Color
	Text          lipgloss.Color
	WarningText   lipgloss.Color
	ErrorText     lipgloss.Color
}

type ThemeColorArgs struct {
	MainBanner    string
	Background    string
	HighlightText string
	Text          string
	WarningText   string
	ErrorText     string
}

func SetThemeColors(args *ThemeColorArgs) *ThemeColors {
	return &ThemeColors{
		MainBanner:    lipgloss.Color(args.MainBanner),
		Background:    lipgloss.Color(args.Background),
		HighlightText: lipgloss.Color(args.HighlightText),
		Text:          lipgloss.Color(args.Text),
		WarningText:   lipgloss.Color(args.WarningText),
		ErrorText:     lipgloss.Color(args.ErrorText),
	}
}
