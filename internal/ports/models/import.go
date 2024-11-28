package models

import (
	"errors"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pimg/certguard/internal/ports/models/commands"
	"github.com/pimg/certguard/internal/ports/models/styles"
)

// keyMap defines a set of keybindings. To work for help it must satisfy
// key.Map. It could also very easily be a map[string]key.Binding.
type importKeyMap struct {
	filepicker.KeyMap
	Back key.Binding
	Quit key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k *importKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Back, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k *importKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Back, k.Quit},
		{k.Up, k.Down},
		{k.Open, k.Select},
	}
}

type clearErrorMsg struct{}

var importKeys = importKeyMap{
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

type ImportModel struct {
	keys         importKeyMap
	styles       *styles.Styles
	filepicker   filepicker.Model
	selectedFile string
	err          error
	commands     *commands.Commands
}

func NewImportModel(cmds *commands.Commands, height int) *ImportModel {
	browseStyle := styles.Theme
	fp := filepicker.New()
	fp.AllowedTypes = []string{".crl", ".pem", ".crt"}
	fp.ShowPermissions = false
	fp.Styles.File = browseStyle.FilePickerFile
	fp.Styles.Selected = browseStyle.FilePickerCurrent
	fp.Styles.Cursor = browseStyle.FilePickerFile
	fp.Height = height - 10
	fp.KeyMap.Back = key.NewBinding(key.WithKeys("h", "backspace", "left"), key.WithHelp("h", "back"))

	fp.CurrentDirectory = cmds.ImportDir()

	return &ImportModel{
		keys:       importKeys,
		styles:     browseStyle,
		filepicker: fp,
		commands:   cmds,
	}
}

func (m *ImportModel) Init() tea.Cmd {
	return m.filepicker.Init()
}

func (m *ImportModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

		cmd := m.commands.ImportFile(m.selectedFile)
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

func (m *ImportModel) View() string {
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
