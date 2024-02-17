package models

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pimg/crl-inspector/internal/ports/models/styles"
)

type sessionState int

const (
	baseView = iota
	inputView
	listView
)

var titles = map[sessionState]string{
	baseView:  "CRL inspector",
	inputView: "Search for a new CRL by entering it's URL",
	listView:  "Pick an entry from the CRL to inspect",
}

type BackMsg int

func Back(state sessionState) tea.Cmd {
	fmt.Printf("Got back command")
	return func() tea.Msg {
		return BackMsg(state)
	}
}

type ExitMsg struct{}

func Exit() tea.Msg {
	return ExitMsg{}
}

// keyMap defines a set of keybindings. To work for help it must satisfy
// key.Map. It could also very easily be a map[string]key.Binding.
type keyMap struct {
	Up     key.Binding
	Down   key.Binding
	Help   key.Binding
	Search key.Binding
	Back   key.Binding
	Quit   key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k *keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k *keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Search},
		{k.Back, k.Help, k.Quit},
	}
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Search: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "search CRL"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back to main view"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

type BaseModel struct {
	title  string
	state  sessionState
	keys   keyMap
	help   help.Model
	styles *styles.Styles
	input  InputModel
}

func NewBaseModel() BaseModel {
	return BaseModel{
		title:  titles[baseView],
		state:  0,
		keys:   keys,
		help:   help.New(),
		styles: styles.DefaultStyles(),
	}
}

func (m BaseModel) Init() tea.Cmd {
	return nil
}

func (m BaseModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd []tea.Cmd

	// global key switches
	switch msg := msg.(type) {
	case BackMsg:
		m.state = sessionState(msg)
		m.title = titles[m.state]
	case ExitMsg:
		return m, tea.Quit
	case tea.WindowSizeMsg:
		m.help.Width = msg.Width
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit) && m.state != inputView: // input view has it's own quit keybinding since we cannot use "q"
			return m, Exit
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		}
	}

	// state specific actions
	switch m.state {
	case inputView:
		switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			m.help.Width = msg.Width
			m.styles.Background.Width(msg.Width)
			m.styles.Background.Height(msg.Height)
		default:
			inputModel, inputCmd := m.input.Update(msg)
			m.input = inputModel.(InputModel)
			cmd = append(cmd, inputCmd)
		}
	case baseView:
		switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			m.help.Width = msg.Width
		case tea.KeyMsg:
			if key.Matches(msg, m.keys.Search) {
				m.state = inputView
				m.title = titles[m.state]
				m.input = NewInputModel()
				return m, m.input.Init()
			}
		}
	}

	return m, tea.Batch(cmd...)
}

func (m BaseModel) View() string {
	helpView := m.help.View(&keys)

	switch m.state {
	case inputView:
		title := m.styles.Title.Render(m.title)
		inputBox := m.input.View()
		helpMenu := m.input.help.View(&inputKeys)
		return lipgloss.JoinVertical(lipgloss.Top, title, inputBox, helpMenu)
	default:
		title := m.styles.Title.Render(m.title)
		helpMenu := helpView
		return lipgloss.JoinVertical(lipgloss.Top, title, helpMenu)
	}
}
