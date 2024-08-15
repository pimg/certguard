package models

import (
	"crypto/x509"
	"testing"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	cmds "github.com/pimg/certguard/internal/ports/models/commands"
	"github.com/pimg/certguard/internal/ports/models/messages"
	"github.com/pimg/certguard/pkg/domain/crl"
	"github.com/stretchr/testify/assert"
)

func TestInitialState(t *testing.T) {
	baseModel := NewBaseModel(nil)

	assert.Equal(t, baseView, baseModel.state)
	assert.Equal(t, titles[baseView], baseModel.title)
	assert.Equal(t, keys, baseModel.keys)
}

func TestSwitchToInputModel(t *testing.T) {
	baseModel := NewBaseModel(nil)

	updatedModel, _ := baseModel.Update(keyBindingToKeyMsg(keys.Download))

	assert.Equal(t, inputView, updatedModel.(BaseModel).state)
	assert.Equal(t, titles[inputView], updatedModel.(BaseModel).title)
}

func TestSwitchBackToBaseModel(t *testing.T) {
	baseModel := NewBaseModel(nil)

	updatedModel, _ := baseModel.Update(keyBindingToKeyMsg(keys.Download))

	assert.Equal(t, inputView, updatedModel.(BaseModel).state)
	assert.Equal(t, titles[inputView], updatedModel.(BaseModel).title)

	updatedModel, _ = updatedModel.Update(keyBindingToKeyMsg(keys.Back))

	assert.Equal(t, baseView, updatedModel.(BaseModel).state)
	assert.Equal(t, titles[baseView], updatedModel.(BaseModel).title)
}

func TestSwitchToBrowseModel(t *testing.T) {
	storage, err := crl.NewMockStorage()
	assert.NoError(t, err)

	commands := cmds.NewCommands(storage)
	baseModel := NewBaseModel(commands)

	updatedModel, _ := baseModel.Update(keyBindingToKeyMsg(keys.Import))

	assert.Equal(t, importView, updatedModel.(BaseModel).state)
	assert.Equal(t, titles[importView], updatedModel.(BaseModel).title)
}

func TestSwitchToListModel(t *testing.T) {
	baseModel := NewBaseModel(nil)

	updatedModel, _ := baseModel.Update(messages.CRLResponseMsg{RevocationList: &x509.RevocationList{}})
	assert.Equal(t, listView, updatedModel.(BaseModel).state)
	assert.Equal(t, titles[listView], updatedModel.(BaseModel).title)
}

func TestSwitchToHome(t *testing.T) {
	baseModel := NewBaseModel(nil)

	updatedModel, _ := baseModel.Update(messages.CRLResponseMsg{RevocationList: &x509.RevocationList{}})
	assert.Equal(t, listView, updatedModel.(BaseModel).state)
	assert.Equal(t, titles[listView], updatedModel.(BaseModel).title)

	updatedModel, _ = updatedModel.Update(keyBindingToKeyMsg(keys.Home))
	assert.Equal(t, baseView, updatedModel.(BaseModel).state)
	assert.Equal(t, titles[baseView], updatedModel.(BaseModel).title)
}

func keyBindingToKeyMsg(keyBinding key.Binding) tea.KeyMsg {
	stringsSlice := keyBinding.Keys()
	var runesSlice []rune
	for _, str := range stringsSlice {
		for _, char := range str {
			runesSlice = append(runesSlice, char)
		}
	}

	return tea.KeyMsg{
		Type:  -1,
		Runes: runesSlice,
		Alt:   false,
	}
}
