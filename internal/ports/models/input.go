package models

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pimg/crl-inspector/internal/ports/models/styles"
)

// keyMap defines a set of keybindings. To work for help it must satisfy
// key.Map. It could also very easily be a map[string]key.Binding.
type inputKeyMap struct {
	Back key.Binding
	Quit key.Binding
	// TODO add enter as confirm
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
		{k.Back, k.Quit},
	}
}

var inputKeys = inputKeyMap{
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back to main view"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	),
}

type InputModel struct {
	keys      inputKeyMap
	help      help.Model
	textinput textinput.Model
	styles    *styles.Styles
}

func NewInputModel() *InputModel {
	i := InputModel{}

	model := textinput.New()
	model.Placeholder = "Enter the URL of a CRL"
	model.Focus()
	i.textinput = model
	i.help = help.New()
	i.keys = inputKeys
	i.styles = styles.DefaultStyles()

	return &i
}

func (i InputModel) Init() tea.Cmd {
	return textinput.Blink
}

func (i InputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, i.keys.Back):
			return i, Back(0)
		case key.Matches(msg, i.keys.Quit):
			return i, Exit // send quitting msg
		} // TODO on "enter" return cmd to perform IO and save input model
	}

	var cmd tea.Cmd
	i.textinput, cmd = i.textinput.Update(msg)
	return i, cmd
}

func (i InputModel) View() string {

	return i.styles.InputField.Render(i.textinput.View())
}
