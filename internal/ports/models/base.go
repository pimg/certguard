package models

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pimg/certguard/internal/ports/models/messages"
	"github.com/pimg/certguard/internal/ports/models/styles"
)

type sessionState int

// TODO consider making a history []sessionState that acts like a stack

const (
	baseView sessionState = iota
	inputView
	listView
	importView
	browseView
	revokedCertificateView
)

var titles = map[sessionState]string{
	baseView:               "CRL inspector",
	inputView:              "Download a new CRL by entering it's URL",
	listView:               "Pick an entry from the CRL to inspect",
	importView:             "Import an existing CRL, from the file system",
	browseView:             "Browse all loaded CRL's from storage",
	revokedCertificateView: "Revoked Certificate",
}

// keyMap defines a set of keybindings. To work for help it must satisfy
// key.Map. It could also very easily be a map[string]key.Binding.
type keyMap struct {
	Help     key.Binding
	Download key.Binding
	Back     key.Binding
	Home     key.Binding
	Import   key.Binding
	Browse   key.Binding
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
		{k.Download, k.Import, k.Home},
		{k.Back, k.Help, k.Quit},
	}
}

var keys = keyMap{
	Download: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "download CRL"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back to previous view"),
	),
	Home: key.NewBinding(
		key.WithKeys("h"),
		key.WithHelp("h", "back to the main view"),
	),
	Import: key.NewBinding(
		key.WithKeys("i"),
		key.WithHelp("i", "import crl from local import directory"),
	),
	Browse: key.NewBinding(
		key.WithKeys("b"),
		key.WithHelp("b", "browse all loaded CRL's from storage"),
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
	title       string
	state       sessionState
	prevState   sessionState
	keys        keyMap
	help        help.Model
	styles      *styles.Styles
	input       InputModel
	browse      *BrowseModel
	list        ListModel
	importModel ImportModel
	err         error
	width       int
	height      int
}

func NewBaseModel() BaseModel {
	return BaseModel{
		title:     titles[baseView],
		state:     baseView,
		prevState: baseView,
		keys:      keys,
		help:      help.New(),
		styles:    styles.DefaultStyles(),
	}
}

func (m BaseModel) Init() tea.Cmd {
	return nil
}

func (m BaseModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		case key.Matches(msg, m.keys.Home) && m.state != inputView:
			m.prevState = m.state
			m.state = baseView
			m.title = titles[baseView]
		case key.Matches(msg, m.keys.Back):
			previousState := m.prevState
			state := m.state
			m.state = previousState
			m.prevState = state
			m.title = titles[m.state]
		}
	case messages.CRLResponseMsg:
		m.prevState = m.state
		m.state = listView
		m.title = titles[listView]
		m.list = NewListModel(msg.RevocationList, msg.URL, m.width, m.height)
	}

	return m.handleStates(msg)
}

// state specific actions
func (m BaseModel) handleStates(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd []tea.Cmd
	switch m.state {
	case inputView:
		inputModel, inputCmd := m.input.Update(msg)
		m.input = inputModel.(InputModel)
		cmd = append(cmd, inputCmd)
	case listView:
		listModel, listCmd := m.list.Update(msg)
		m.list = listModel.(ListModel)

		if m.list.selectedItem != nil && m.list.itemSelected {
			m.prevState = m.state
			m.state = revokedCertificateView
			m.title = titles[revokedCertificateView]
		}
		cmd = append(cmd, listCmd)
	case revokedCertificateView:
		revokedCertificateModel, revokedCertificateCmd := m.list.selectedItem.Update(msg)
		rcm := revokedCertificateModel.(RevokedCertificateModel)
		m.list.selectedItem = &rcm
		cmd = append(cmd, revokedCertificateCmd)
	case importView:
		importModel, importCmd := m.importModel.Update(msg)
		m.importModel = importModel.(ImportModel)
		cmd = append(cmd, importCmd)
	case browseView:
		browseModel, browseCmd := m.browse.Update(msg)
		m.browse = browseModel.(*BrowseModel)
		cmd = append(cmd, browseCmd)
	case baseView:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if key.Matches(msg, m.keys.Download) {
				m.prevState = m.state
				m.state = inputView
				m.title = titles[m.state]
				m.input = NewInputModel()
				return m, m.input.Init()
			}
			if key.Matches(msg, m.keys.Import) {
				m.prevState = m.state
				m.state = importView
				m.title = titles[m.state]
				m.importModel = NewImportModel()
				return m, m.importModel.Init()
			}
			if key.Matches(msg, m.keys.Browse) {
				m.prevState = m.state
				m.state = browseView
				m.title = titles[m.state]
				m.browse = NewBrowseModel(m.height)
				return m, m.browse.Init()
			}
		}
	}

	return m, tea.Batch(cmd...)
}

func (m BaseModel) View() string {
	errorMsg := ""

	switch m.state {
	case inputView:
		title := m.styles.Title.Render(m.title)
		inputBox := m.input.View()
		helpMenu := m.help.View(&inputKeys)
		height := strings.Count(inputBox, "\n") + strings.Count(title, "\n")
		return lipgloss.JoinVertical(lipgloss.Top, title, inputBox) + lipgloss.Place(m.width, m.height-height-1, lipgloss.Left, lipgloss.Bottom, helpMenu)
	case listView:
		title := m.styles.Title.Render(m.title)
		listInfo := m.list.View()
		return lipgloss.JoinVertical(lipgloss.Top, title, listInfo)
	case revokedCertificateView:
		title := m.styles.Title.Render(m.title)
		helpMenu := m.help.View(&revokedCertificateKeys)
		revokedCertificateDetails := m.list.selectedItem.View()
		height := strings.Count(revokedCertificateDetails, "\n") + strings.Count(title, "\n")
		return lipgloss.JoinVertical(lipgloss.Top, title, revokedCertificateDetails) + lipgloss.Place(m.width, m.height-height-1, lipgloss.Left, lipgloss.Bottom, helpMenu)
	case importView:
		title := m.styles.Title.Render(m.title)
		listInfo := m.importModel.View()
		helpMenu := m.help.View(&listKeys)
		height := strings.Count(listInfo, "\n") + strings.Count(title, "\n")
		return lipgloss.JoinVertical(lipgloss.Top, title, listInfo) + lipgloss.Place(m.width, m.height-height-1, lipgloss.Left, lipgloss.Bottom, helpMenu)
	case browseView:
		title := m.styles.Title.Render(m.title)
		table := m.browse.View()
		helpMenu := m.help.View(&browseKeys)
		height := strings.Count(table, "\n") + strings.Count(title, "\n")
		return lipgloss.JoinVertical(lipgloss.Top, title, table) + lipgloss.Place(m.width, m.height-height-1, lipgloss.Left, lipgloss.Bottom, helpMenu)
	default:
		title := m.styles.Title.Render(m.title)
		if m.err != nil {
			errorMsg = m.err.Error()
		}
		downloadHelp := m.styles.BaseText.Render("Download a CRL file: ") + "d"
		importHelp := m.styles.BaseText.Render("Import a CRL file from local import directory: ") + "i"
		browseHelp := m.styles.BaseText.Render("Browse all loaded CRL's from storage") + "b"
		mainMenu := fmt.Sprintf("%s\n%s\n%s", downloadHelp, importHelp, browseHelp)
		helpMenu := m.help.View(&keys)
		height := strings.Count(title, "\n")
		return lipgloss.JoinVertical(lipgloss.Top, title, errorMsg, mainMenu) + lipgloss.Place(m.width, m.height-height-5, lipgloss.Left, lipgloss.Bottom, helpMenu)
	}
}
