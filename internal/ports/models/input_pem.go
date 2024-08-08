package models

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pimg/certguard/internal/ports/models/commands"
	"github.com/pimg/certguard/internal/ports/models/messages"
	"github.com/pimg/certguard/internal/ports/models/styles"
)

// keyMap defines a set of keybindings. To work for help it must satisfy
// key.Map. It could also very easily be a map[string]key.Binding.
type inputPemKeyMap struct {
	Back  key.Binding
	Enter key.Binding
	Quit  key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k *inputPemKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Back, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k *inputPemKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Back, k.Enter, k.Quit},
	}
}

var inputPemKeys = inputPemKeyMap{
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back to main view"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithKeys("enter", "submit PEM")),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	),
}

type InputPemModel struct {
	keys     inputPemKeyMap
	styles   *styles.Styles
	textArea textarea.Model
	msg      string // Tmp placeholder before we implement the view model
}

func NewInputPemModel(height, width int) *InputPemModel {
	i := &InputPemModel{
		keys:   inputPemKeys,
		styles: styles.DefaultStyles(),
	}

	ta := textarea.New()
	ta.Placeholder = "Insert PEM content here..."
	ta.CharLimit = 0
	ta.Focus()

	if width > 80 {
		ta.SetWidth(80)
	} else {
		ta.SetWidth(width)
	}

	ta.SetHeight(height / 2)

	i.textArea = ta

	return i
}

func (i *InputPemModel) Init() tea.Cmd {
	return textarea.Blink
}

func (i *InputPemModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, i.keys.Quit):
			return i, tea.Quit
		case key.Matches(msg, i.keys.Enter):
			cmd = commands.ParsePemCertficate(i.textArea.Value())
			i.textArea.Reset()
			return i, cmd
		default:
			if !i.textArea.Focused() {
				cmd = i.textArea.Focus()
				cmds = append(cmds, cmd)
			}
		}
	case messages.ErrorMsg:
		i.msg = msg.Err.Error()
	}

	i.textArea, cmd = i.textArea.Update(msg)
	cmds = append(cmds, cmd)
	return i, tea.Batch(cmds...)
}

func (i *InputPemModel) View() string {
	var s strings.Builder

	if i.textArea.Err != nil {
		s.WriteString(i.styles.ErrorMessages.Render(i.textArea.Err.Error()))
	}

	if i.msg != "" {
		s.WriteString(i.styles.BaseText.Render(i.msg))
	}

	return lipgloss.JoinVertical(lipgloss.Top, i.styles.TextArea.Render(i.textArea.View()), s.String())
}
