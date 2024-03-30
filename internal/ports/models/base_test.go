package models

import (
	"crypto/x509"
	"testing"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pimg/certguard/internal/adapter"
	"github.com/pimg/certguard/internal/ports/models/messages"
	"github.com/stretchr/testify/assert"
)

func TestInitialState(t *testing.T) {
	baseModel := NewBaseModel()

	assert.Equal(t, baseView, baseModel.state)
	assert.Equal(t, titles[baseView], baseModel.title)
	assert.Equal(t, keys, baseModel.keys)
}

func TestSwitchToInputModel(t *testing.T) {
	baseModel := NewBaseModel()

	updatedModel, _ := baseModel.Update(keyBindingToKeyMsg(keys.Download))

	assert.Equal(t, inputView, updatedModel.(BaseModel).state)
	assert.Equal(t, titles[inputView], updatedModel.(BaseModel).title)
}

func TestSwitchBackToBaseModel(t *testing.T) {
	baseModel := NewBaseModel()

	updatedModel, _ := baseModel.Update(keyBindingToKeyMsg(keys.Download))

	assert.Equal(t, inputView, updatedModel.(BaseModel).state)
	assert.Equal(t, titles[inputView], updatedModel.(BaseModel).title)

	updatedModel, _ = updatedModel.Update(keyBindingToKeyMsg(keys.Back))

	assert.Equal(t, baseView, updatedModel.(BaseModel).state)
	assert.Equal(t, titles[baseView], updatedModel.(BaseModel).title)
}

func TestSwitchToBrowseModel(t *testing.T) {
	_, err := adapter.NewFileCache()
	assert.NoError(t, err)
	baseModel := NewBaseModel()

	updatedModel, _ := baseModel.Update(keyBindingToKeyMsg(keys.Browse))

	assert.Equal(t, browseView, updatedModel.(BaseModel).state)
	assert.Equal(t, titles[browseView], updatedModel.(BaseModel).title)
}

func TestSwitchToListModel(t *testing.T) {
	baseModel := NewBaseModel()

	updatedModel, _ := baseModel.Update(messages.CRLResponseMsg{RevocationList: &x509.RevocationList{}})
	assert.Equal(t, listView, updatedModel.(BaseModel).state)
	assert.Equal(t, titles[listView], updatedModel.(BaseModel).title)
}

func TestSwitchToHome(t *testing.T) {
	baseModel := NewBaseModel()

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
