package models

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pimg/certguard/internal/ports/models/commands"
	"github.com/pimg/certguard/internal/ports/models/messages"
	"github.com/pimg/certguard/internal/ports/models/styles"
	"github.com/pimg/certguard/pkg/uri"
)

// keyMap defines a set of keybindings. To work for help it must satisfy
// key.Map. It could also very easily be a map[string]key.Binding.
type inputKeyMap struct {
	Back  key.Binding
	Enter key.Binding
	Quit  key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k *inputKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Back, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k *inputKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Back, k.Enter, k.Quit},
	}
}

var inputKeys = inputKeyMap{
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back to main view"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithKeys("enter", "confirm to get CRL")),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	),
}

type InputModel struct {
	keys      inputKeyMap
	textinput textinput.Model
	styles    *styles.Styles
}

func NewInputModel() InputModel {
	i := InputModel{}

	input := textinput.New()
	input.Placeholder = "Enter the URL of a CRL"
	input.Focus()
	i.textinput = input
	i.keys = inputKeys
	i.styles = styles.DefaultStyles()

	return i
}

func (i InputModel) Init() tea.Cmd {
	return textinput.Blink
}

func (i InputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		i.textinput.Err = nil
		switch {
		case key.Matches(msg, i.keys.Quit):
			return i, tea.Quit
		case key.Matches(msg, i.keys.Enter):
			confirmedInput := i.textinput.Value()
			i.textinput.Reset()
			url, err := uri.ValidateURI(confirmedInput)
			if err != nil {
				i.textinput.Err = err
				return i, nil
			}
			cmd = commands.GetCRL(url)
			return i, cmd
		}

	case messages.ErrorMsg:
		i.textinput.Err = msg.Err
		return i, cmd
	}

	i.textinput, cmd = i.textinput.Update(msg)
	return i, cmd
}

func (i InputModel) View() string {
	errorMsg := ""
	// TODO introduce spinner since download time can be long toggleling spinners seems to  be done by creating  a new one: https://github.com/charmbracelet/bubbletea/blob/master/examples/spinners/main.go
	if i.textinput.Err != nil {
		errorMsg = i.textinput.Err.Error()
		return lipgloss.JoinVertical(lipgloss.Top, i.styles.InputField.Render(i.textinput.View()), i.styles.ErrorMessages.Render(errorMsg))
	}

	return lipgloss.JoinVertical(lipgloss.Top, i.styles.InputField.Render(i.textinput.View()), errorMsg)
}
