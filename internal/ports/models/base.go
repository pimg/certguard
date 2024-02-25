package models

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pimg/certguard/internal/ports/models/messages"
	"github.com/pimg/certguard/internal/ports/models/styles"
)

type sessionState int

const (
	baseView = iota
	inputView
	listView
)

var titles = map[sessionState]string{
	baseView:  "CRL inspector",
	inputView: "Download a new CRL by entering it's URL",
	listView:  "Pick an entry from the CRL to inspect",
}

// keyMap defines a set of keybindings. To work for help it must satisfy
// key.Map. It could also very easily be a map[string]key.Binding.
type keyMap struct {
	Up       key.Binding
	Down     key.Binding
	Help     key.Binding
	Download key.Binding
	Back     key.Binding
	Quit     key.Binding
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
		{k.Up, k.Down, k.Download},
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
	Download: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "download CRL"),
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
	list   ListModel
	err    error
	width  int
	height int
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
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.styles.Background.Width(msg.Width)
		m.styles.Background.Height(msg.Height)
		m.help.Width = msg.Width
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit) && m.state != inputView: // input view has it's own quit keybinding since we cannot use "q"
			return m, tea.Quit
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.keys.Back):
			m.state = baseView
			m.title = titles[baseView]
		}
	case messages.CRLResponseMsg:
		m.state = listView
		m.title = titles[listView]
		m.list = NewListModel(msg.RevocationList)
	}

	// state specific actions
	switch m.state {
	case inputView:
		switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			m.help.Width = msg.Width
		default:
			inputModel, inputCmd := m.input.Update(msg)
			m.input = inputModel.(InputModel)
			cmd = append(cmd, inputCmd)
		}
	case listView:
		switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			m.help.Width = msg.Width
		default:
			listModel, listCmd := m.list.Update(msg)
			m.list = listModel.(ListModel)
			cmd = append(cmd, listCmd)
		}
	case baseView:
		switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			m.help.Width = msg.Width
		case tea.KeyMsg:
			if key.Matches(msg, m.keys.Download) {
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
	errorMsg := ""

	switch m.state {
	case inputView:
		title := m.styles.Title.Render(m.title)
		inputBox := m.input.View()
		helpMenu := m.input.help.View(&inputKeys)
		height := strings.Count(inputBox, "\n") + strings.Count(title, "\n")
		return lipgloss.JoinVertical(lipgloss.Top, title, inputBox) + lipgloss.Place(m.width, m.height-height, lipgloss.Left, lipgloss.Bottom, helpMenu)
	case listView:
		title := m.styles.Title.Render(m.title)
		listInfo := m.list.View()
		helpMenu := m.list.help.View(&listKeys)
		height := strings.Count(listInfo, "\n") + strings.Count(title, "\n")
		return lipgloss.JoinVertical(lipgloss.Top, title, listInfo) + lipgloss.Place(m.width, m.height-height, lipgloss.Left, lipgloss.Bottom, helpMenu)
	default:
		title := m.styles.Title.Render(m.title)
		if m.err != nil {
			errorMsg = m.err.Error()
		}
		helpMenu := helpView
		height := strings.Count(title, "\n")
		return lipgloss.JoinVertical(lipgloss.Top, title, errorMsg) + lipgloss.Place(m.width, m.height-height-1, lipgloss.Left, lipgloss.Bottom, helpMenu)
	}
}
