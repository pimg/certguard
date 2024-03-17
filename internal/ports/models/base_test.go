package models

import (
	"testing"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
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
