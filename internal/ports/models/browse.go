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
	Back  key.Binding
	Quit  key.Binding
	Enter key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k *browseKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Back, k.Quit, k.LineUp, k.LineDown, k.Enter}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k *browseKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Back, k.Quit},
		{k.LineUp, k.LineDown},
		{k.GotoTop, k.GotoBottom},
		{k.Enter},
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
	KeyMap: table.DefaultKeyMap(),
}

type BrowseModel struct {
	table table.Model
}

func NewBrowseModel(height int) *BrowseModel {
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
		BorderForeground(styles.DefaultStyles().ListComponentTitle).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(styles.DefaultStyles().FilePickerCurrent.GetForeground()).
		Background(styles.DefaultStyles().BaseText.GetBackground()).
		Bold(false)
	tbl.SetStyles(s)

	return &BrowseModel{
		table: tbl,
	}
}

func (m *BrowseModel) Init() tea.Cmd {
	return commands.GetCRLsFromStore
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
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			cmd := commands.GetRevokedCertificates(m.table.SelectedRow()[0], m.table.SelectedRow()[1], m.table.SelectedRow()[2], m.table.SelectedRow()[3])
			return m, cmd
		}
	}
	m.table, cmd = m.table.Update(msg)

	return m, cmd
}

func (m *BrowseModel) View() string {
	var s strings.Builder
	s.WriteString("\n\n" + m.table.View() + "\n")
	return s.String()
}
