package models

import (
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pimg/certguard/internal/ports/models/commands"
	"github.com/pimg/certguard/internal/ports/models/messages"
	"github.com/pimg/certguard/internal/ports/models/styles"
)

// keyMap defines a set of keybindings. To work for help it must satisfy
// key.Map. It could also very easily be a map[string]key.Binding.
type browseKeyMap struct {
	table.KeyMap
	Back   key.Binding
	Quit   key.Binding
	Enter  key.Binding
	Delete key.Binding
	Y      key.Binding
	N      key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k *browseKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Back, k.Quit, k.LineUp, k.LineDown, k.Enter, k.Delete}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k *browseKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Back, k.Quit},
		{k.LineUp, k.LineDown},
		{k.GotoTop, k.GotoBottom},
		{k.Enter, k.Delete},
	}
}

var browseKeys = browseKeyMap{
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back to main view"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select row"),
	),
	Delete: key.NewBinding(
		key.WithKeys("delete"),
		key.WithHelp("delete", "marks a CRL for deletion"),
	),
	Y: key.NewBinding(
		key.WithKeys("y"),
		key.WithHelp("y", "confirm deletion"),
	),
	N: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "cancel deletion"),
	),
	KeyMap: table.DefaultKeyMap(),
}

type BrowseModel struct {
	table             table.Model
	markedForDeletion string
	errorMsg          string
	styles            *styles.Styles
	commands          *commands.Commands
}

func NewBrowseModel(height int, cmds *commands.Commands) *BrowseModel {
	columns := []table.Column{
		{Title: "ID", Width: 2},
		{Title: "Name", Width: 30},
		{Title: "This Update", Width: 11},
		{Title: "Next Update", Width: 11},
		{Title: "Url", Width: 15},
	}

	tbl := table.New(table.WithColumns(columns), table.WithFocused(true), table.WithHeight(height-10), table.WithWidth(80))
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(styles.Theme.ListComponentTitle).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(styles.Theme.FilePickerCurrent.GetForeground()).
		Background(styles.Theme.BaseText.GetBackground()).
		Bold(false)
	tbl.SetStyles(s)

	return &BrowseModel{
		table:    tbl,
		styles:   styles.Theme,
		commands: cmds,
	}
}

func (m *BrowseModel) Init() tea.Cmd {
	return m.commands.GetCRLsFromStore
}

func (m *BrowseModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case messages.ListCRLsResponseMsg:
		rows := make([]table.Row, len(msg.CRLs))
		for i, CRL := range msg.CRLs {
			rows[i] = table.Row{
				strconv.Itoa(int(CRL.ID)),
				CRL.Name,
				CRL.ThisUpdate.Format(time.DateOnly),
				CRL.NextUpdate.Format(time.DateOnly),
				CRL.URL.String(),
			}
		}
		m.table.SetRows(rows)
	case messages.CRLDeleteConfirmationMsg:
		if msg.DeletionSuccessful {
			m.deleteFromRows()
		}
	case messages.ErrorMsg:
		m.errorMsg = msg.Err.Error()
		m.markedForDeletion = ""
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			cmd := m.commands.GetRevokedCertificates(&commands.GetRevokedCertificatesArgs{
				ID:         m.table.SelectedRow()[0],
				CN:         m.table.SelectedRow()[1],
				ThisUpdate: m.table.SelectedRow()[2],
				NextUpdate: m.table.SelectedRow()[3],
				URL:        m.table.SelectedRow()[4],
			})
			return m, cmd
		case "delete":
			m.markedForDeletion = m.table.SelectedRow()[0]
		case "n":
			m.markedForDeletion = ""
		case "y":
			if m.markedForDeletion != "" {
				cmd := m.commands.DeleteCRLFromStore(m.markedForDeletion)
				return m, cmd
			}
		}
	}
	m.table, cmd = m.table.Update(msg)

	return m, cmd
}

func (m *BrowseModel) deleteFromRows() {
	rows := m.table.Rows()
	for i, row := range rows {
		if row[0] == m.markedForDeletion {
			m.table.SetRows(append(rows[:i], rows[i+1:]...))
		}
	}
	m.markedForDeletion = ""
}

func (m *BrowseModel) View() string {
	var s strings.Builder

	if m.markedForDeletion != "" {
		s.WriteString(m.styles.WarningText.Render("\n\n Do you want to delete CRL : " + m.markedForDeletion + " y(es) n(o)"))
	}

	if m.errorMsg != "" {
		s.WriteString(m.styles.WarningText.Render("\n\n" + m.errorMsg))
	}

	s.WriteString("\n\n" + m.table.View())
	return s.String()
}
