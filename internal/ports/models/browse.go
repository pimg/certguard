package models

import (
	"errors"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pimg/certguard/internal/adapter"
	"github.com/pimg/certguard/internal/ports/models/commands"
	"github.com/pimg/certguard/internal/ports/models/styles"
)

// keyMap defines a set of keybindings. To work for help it must satisfy
// key.Map. It could also very easily be a map[string]key.Binding.
type browseKeyMap struct {
	filepicker.KeyMap
	Back key.Binding
	Quit key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k *browseKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Back, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k *browseKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Back, k.Quit},
		{k.Up, k.Down},
		{k.Open, k.Select},
	}
}

type clearErrorMsg struct{}

var browseKeys = browseKeyMap{
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back to main view"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	KeyMap: filepicker.DefaultKeyMap(),
}

type BrowseModel struct {
	keys         browseKeyMap
	styles       *styles.Styles
	filepicker   filepicker.Model
	selectedFile string
	err          error
}

func NewBrowseModel() BrowseModel {
	browseStyle := styles.DefaultStyles()
	fp := filepicker.New()
	fp.AllowedTypes = []string{".crl", ".pem", ".crt"}
	fp.ShowPermissions = false
	fp.Styles.File = browseStyle.FilePickerFile
	fp.Styles.Selected = browseStyle.FilePickerCurrent
	fp.Styles.Cursor = browseStyle.FilePickerFile
	fp.KeyMap.Back = key.NewBinding(key.WithKeys("h", "backspace", "left"), key.WithHelp("h", "back"))

	fp.CurrentDirectory = adapter.GlobalCache.(*adapter.FileCache).Dir()

	return BrowseModel{
		keys:       browseKeys,
		styles:     styles.DefaultStyles(),
		filepicker: fp,
	}
}

func (m BrowseModel) Init() tea.Cmd {
	return m.filepicker.Init()
}

func (m BrowseModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case clearErrorMsg:
		m.err = nil
	}

	var cmd tea.Cmd
	m.filepicker, cmd = m.filepicker.Update(msg)

	// Did the user select a file?
	if didSelect, path := m.filepicker.DidSelectFile(msg); didSelect {
		// Get the path of the selected file.
		m.selectedFile = path

		cmd := commands.LoadCRL(m.selectedFile)
		m.selectedFile = ""
		m.filepicker.Path = ""
		return m, cmd
	}

	// Did the user select a disabled file?
	// This is only necessary to display an error to the user.
	if didSelect, path := m.filepicker.DidSelectDisabledFile(msg); didSelect {
		// Let's clear the selectedFile and display an error.
		m.err = errors.New(path + " is not valid.")
		m.selectedFile = ""
		return m, tea.Batch(cmd, clearErrorAfter(2*time.Second))
	}

	return m, cmd
}

func (m BrowseModel) View() string {
	var s strings.Builder
	s.WriteString("\n  ")
	if m.err != nil {
		s.WriteString(m.filepicker.Styles.DisabledFile.Render(m.err.Error()))
	} else {
		s.WriteString("Pick a file:")
	}
	s.WriteString("\n\n" + m.filepicker.View() + "\n")
	return s.String()
}

func clearErrorAfter(t time.Duration) tea.Cmd {
	return tea.Tick(t, func(_ time.Time) tea.Msg {
		return clearErrorMsg{}
	})
}
