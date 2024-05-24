package models

import (
	"crypto/x509"
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pimg/certguard/internal/ports/models/styles"
)

// keyMap defines a set of keybindings. To work for help it must satisfy
// key.Map. It could also very easily be a map[string]key.Binding.
type listKeyMap struct {
	list.KeyMap
	Back   key.Binding
	Quit   key.Binding
	Select key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k *listKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Back, k.Quit, k.Select}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k *listKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Back, k.Quit},
	}
}

var listKeys = listKeyMap{
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back to main view"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Select: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
}

type item struct {
	serialnumber, revocationReason, revocationDate string
}

func (i item) Title() string       { return i.serialnumber }
func (i item) Description() string { return i.revocationDate }
func (i item) FilterValue() string { return i.serialnumber }

const TOP_INFO_HEIGHT = 12

type ListModel struct {
	keys         listKeyMap
	styles       *styles.Styles
	list         list.Model
	crl          *x509.RevocationList
	selectedItem *RevokedCertificateModel
	itemSelected bool
}

func NewListModel(crl *x509.RevocationList, width, height int) ListModel {
	items := revokedCertificatesToItems(crl.RevokedCertificateEntries)

	defaultDelegate := list.NewDefaultDelegate()
	c := styles.DefaultStyles().ListComponentTitle
	defaultDelegate.Styles.SelectedTitle = defaultDelegate.Styles.SelectedTitle.Foreground(c).BorderLeftForeground(c)
	defaultDelegate.Styles.SelectedDesc = defaultDelegate.Styles.SelectedTitle

	revokedList := list.New(items, defaultDelegate, width, height-TOP_INFO_HEIGHT)
	revokedList.Title = "Revoked Certificates"
	revokedList.KeyMap.Quit = key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("ctrl-c", "quit"))
	revokedList.KeyMap.ClearFilter = key.NewBinding(key.WithKeys("ctrl+q"), key.WithHelp("ctrl-q", "clear"))
	revokedList.KeyMap.CancelWhileFiltering = key.NewBinding(key.WithKeys("ctrl+q"), key.WithHelp("ctrl-q", "clear"))
	revokedList.AdditionalShortHelpKeys = listKeys.ShortHelp

	revokedList.Styles.Title = revokedList.Styles.Title.Background(c)
	return ListModel{
		keys:   listKeys,
		styles: styles.DefaultStyles(),
		list:   revokedList,
		crl:    crl,
	}
}

func revokedCertificatesToItems(entries []x509.RevocationListEntry) []list.Item {
	items := make([]list.Item, 0, len(entries))
	for _, entry := range entries {
		items = append(items, item{
			serialnumber:     entry.SerialNumber.String(),
			revocationReason: strconv.Itoa(entry.ReasonCode),
			revocationDate:   entry.RevocationTime.String(),
		})
	}

	return items
}

func (l ListModel) Init() tea.Cmd {
	return nil
}

func (l ListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, listKeys.Select):
			if len(l.list.VisibleItems()) != 0 {
				selectedItem := l.list.SelectedItem().(item)
				revokedCertificateModel := NewRevokedCertificateModel(selectedItem.serialnumber, selectedItem.revocationReason, selectedItem.revocationDate)
				l.selectedItem = &revokedCertificateModel
				l.itemSelected = true
			}
		default:
			l.itemSelected = false
		}
	case tea.WindowSizeMsg:
		l.list.SetSize(msg.Width, msg.Height-TOP_INFO_HEIGHT)
	}

	l.list, cmd = l.list.Update(msg)
	return l, cmd
}

func (l ListModel) View() string {
	issuer := fmt.Sprintf("CRL Issuer          : %s", l.crl.Issuer)
	updatedAt := fmt.Sprintf("Updated At          : %s", l.crl.ThisUpdate)
	nextUpdate := fmt.Sprintf("Next Update         : %s", l.crl.NextUpdate)
	revokedCertCount := fmt.Sprintf("Revoked Certificates: %d", len(l.crl.RevokedCertificateEntries))
	revokedList := l.list.View()

	crlInfo := l.styles.Text.Render(
		fmt.Sprintf("%s\n%s\n%s\n%s", issuer, updatedAt, nextUpdate, revokedCertCount),
	)

	return lipgloss.JoinVertical(lipgloss.Top, crlInfo, revokedList)
}
