package models

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pimg/certguard/internal/ports/models/commands"
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
	importPemView
	inputPemView
	certificateView
)

var titles = map[sessionState]string{
	baseView:               "CRL inspector",
	inputView:              "Download a new CRL by entering it's URL",
	listView:               "Pick an entry from the CRL to inspect",
	importView:             "Import an existing CRL, from the file system",
	browseView:             "Browse all loaded CRL's from storage",
	revokedCertificateView: "Revoked Certificate",
	importPemView:          "Import a PEM certificate",
	inputPemView:           "Input a PEM certificate",
	certificateView:        "view a parsed certificate",
}

// keyMap defines a set of keybindings. To work for help it must satisfy
// key.Map. It could also very easily be a map[string]key.Binding.
type keyMap struct {
	Help      key.Binding
	Download  key.Binding
	Back      key.Binding
	Home      key.Binding
	Import    key.Binding
	Browse    key.Binding
	InputPem  key.Binding
	ImportPem key.Binding
	Quit      key.Binding
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
		key.WithHelp("b", "browseModel all loaded CRL's from storage"),
	),
	InputPem: key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "inputModel a PEM certificate"),
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
	title            string
	state            sessionState
	prevState        sessionState
	keys             keyMap
	help             help.Model
	styles           *styles.Styles
	commands         *commands.Commands
	inputModel       *InputModel
	browseModel      *BrowseModel
	listModel        *ListModel
	importModel      *ImportModel
	inputPemModel    *InputPemModel
	certificateModel *CertificateModel
	err              error
	width            int
	height           int
}

func NewBaseModel(cmds *commands.Commands) BaseModel {
	return BaseModel{
		title:     titles[baseView],
		state:     baseView,
		prevState: baseView,
		keys:      keys,
		help:      help.New(),
		styles:    styles.DefaultStyles(),
		commands:  cmds,
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
		case key.Matches(msg, m.keys.Quit) && m.state != inputView && m.state != inputPemView: // inputModel view has it's own quit keybinding since we cannot use "q"
			return m, tea.Quit
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.keys.Home) && m.state != inputView && m.state != inputPemView:
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
		m.listModel = NewListModel(msg.RevocationList, msg.URL, m.width, m.height, m.commands)
	case messages.PemCertificateMsg:
		m.prevState = m.state
		m.state = certificateView
		m.title = titles[certificateView]
		m.certificateModel = NewCertificateModel(msg.Certificate, msg.CertificateChain, m.commands)
	}

	return m.handleStates(msg)
}

// state specific actions
func (m BaseModel) handleStates(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd []tea.Cmd
	switch m.state {
	case inputView:
		inputModel, inputCmd := m.inputModel.Update(msg)
		m.inputModel = inputModel.(*InputModel)
		cmd = append(cmd, inputCmd)
	case listView:
		listModel, listCmd := m.listModel.Update(msg)
		m.listModel = listModel.(*ListModel)

		if m.listModel.selectedItem != nil && m.listModel.itemSelected {
			m.prevState = m.state
			m.state = revokedCertificateView
			m.title = titles[revokedCertificateView]
		}
		cmd = append(cmd, listCmd)
	case revokedCertificateView:
		revokedCertificateModel, revokedCertificateCmd := m.listModel.selectedItem.Update(msg)
		rcm := revokedCertificateModel.(*RevokedCertificateModel)
		m.listModel.selectedItem = rcm
		cmd = append(cmd, revokedCertificateCmd)
	case importView:
		importModel, importCmd := m.importModel.Update(msg)
		m.importModel = importModel.(*ImportModel)
		cmd = append(cmd, importCmd)
	case browseView:
		browseModel, browseCmd := m.browseModel.Update(msg)
		m.browseModel = browseModel.(*BrowseModel)
		cmd = append(cmd, browseCmd)
	case inputPemView:
		inputPemModel, inputPemCmd := m.inputPemModel.Update(msg)
		m.inputPemModel = inputPemModel.(*InputPemModel)
		cmd = append(cmd, inputPemCmd)
	case certificateView:
		certificateModel, certificateCmd := m.certificateModel.Update(msg)
		m.certificateModel = certificateModel.(*CertificateModel)
		cmd = append(cmd, certificateCmd)
	case baseView:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if key.Matches(msg, m.keys.Download) {
				m.prevState = m.state
				m.state = inputView
				m.title = titles[m.state]
				m.inputModel = NewInputModel(m.commands)
				return m, m.inputModel.Init()
			}
			if key.Matches(msg, m.keys.Import) {
				m.prevState = m.state
				m.state = importView
				m.title = titles[m.state]
				m.importModel = NewImportModel(m.commands, m.height)
				return m, m.importModel.Init()
			}
			if key.Matches(msg, m.keys.Browse) {
				m.prevState = m.state
				m.state = browseView
				m.title = titles[m.state]
				m.browseModel = NewBrowseModel(m.height, m.commands)
				return m, m.browseModel.Init()
			}
			if key.Matches(msg, m.keys.InputPem) {
				m.prevState = m.state
				m.state = inputPemView
				m.title = titles[m.state]
				m.inputPemModel = NewInputPemModel(m.height, m.width, m.commands)
				return m, m.inputModel.Init()
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
		inputBox := m.inputModel.View()
		helpMenu := m.help.View(&inputKeys)
		height := strings.Count(inputBox, "\n") + strings.Count(title, "\n")
		return lipgloss.JoinVertical(lipgloss.Top, title, inputBox) + lipgloss.Place(m.width, m.height-height-1, lipgloss.Left, lipgloss.Bottom, helpMenu)
	case listView:
		title := m.styles.Title.Render(m.title)
		listInfo := m.listModel.View()
		return lipgloss.JoinVertical(lipgloss.Top, title, listInfo)
	case revokedCertificateView:
		title := m.styles.Title.Render(m.title)
		helpMenu := m.help.View(&revokedCertificateKeys)
		revokedCertificateDetails := m.listModel.selectedItem.View()
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
		table := m.browseModel.View()
		helpMenu := m.help.View(&browseKeys)
		height := strings.Count(table, "\n") + strings.Count(title, "\n")
		return lipgloss.JoinVertical(lipgloss.Top, title, table) + lipgloss.Place(m.width, m.height-height-1, lipgloss.Left, lipgloss.Bottom, helpMenu)
	case inputPemView:
		title := m.styles.Title.Render(m.title)
		textArea := m.inputPemModel.View()
		helpMenu := m.help.View(&inputPemKeys)
		height := strings.Count(textArea, "\n") + strings.Count(title, "\n")
		return lipgloss.JoinVertical(lipgloss.Top, title, textArea) + lipgloss.Place(m.width, m.height-height-1, lipgloss.Left, lipgloss.Bottom, helpMenu)
	case certificateView:
		title := m.styles.Title.Render(m.title)
		certInfo := m.certificateModel.View()
		helpMenu := m.help.View(&certificateKeys)
		height := strings.Count(certInfo, "\n") + strings.Count(title, "\n")
		return lipgloss.JoinVertical(lipgloss.Top, title, certInfo) + lipgloss.Place(m.width, m.height-height-1, lipgloss.Left, lipgloss.Bottom, helpMenu)
	default:
		title := m.styles.Title.Render(m.title)
		if m.err != nil {
			errorMsg = m.err.Error()
		}
		downloadHelp := m.styles.BaseMenuText.Render("Download a CRL file: ") + "d"
		importHelp := m.styles.BaseMenuText.Render("Import a CRL, or PEM Cert from import directory: ") + "i"
		browseHelp := m.styles.BaseMenuText.Render("Browse all loaded CRL's from storage") + "b"
		mainMenu := fmt.Sprintf("%s\n%s\n%s", downloadHelp, importHelp, browseHelp)

		inputPemHelp := m.styles.BaseMenuText.Render("Input a Certificate in PEM format") + "p"
		pemMenu := fmt.Sprintf("%s\n", inputPemHelp)

		menu := fmt.Sprintf("%s\n\n%s", mainMenu, pemMenu)

		helpMenu := m.help.View(&keys)
		height := strings.Count(title, "\n")
		return lipgloss.JoinVertical(lipgloss.Top, title, errorMsg, menu) + lipgloss.Place(m.width, m.height-height-7, lipgloss.Left, lipgloss.Bottom, helpMenu)
	}
}
