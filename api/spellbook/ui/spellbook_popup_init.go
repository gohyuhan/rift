package ui

import (
	"charm.land/bubbles/v2/textinput"
	"github.com/gohyuhan/rift/api/learn"
	"github.com/gohyuhan/rift/i18n"
)

// ----------------------------------
//
//	initialises the LearnPopUpModel with two text inputs (spell name and spell
//	command), configures placeholders from i18n, focuses the spell name input,
//	and stores the populated model on the parent SpellbookInteractiveModel;
//	OnInputFuncTrigger is set to learn.SaveSpell so submissions are persisted
//	to the database
//
// ----------------------------------
func initLearnPopUpModel(m *SpellbookInteractiveModel) {
	spellNameInput := textinput.New()
	spellNameInput.SetValue("")
	spellNameInput.Placeholder = i18n.LANGUAGEMAPPING.SpellNameInputPlaceHolder
	spellNameInput.Focus()
	spellNameInput.SetVirtualCursor(true)

	spellCmdInput := textinput.New()
	spellCmdInput.SetValue("")
	spellCmdInput.Placeholder = i18n.LANGUAGEMAPPING.SpellCommandInputPlaceHolder
	spellCmdInput.SetVirtualCursor(true)

	popUpModel := &LearnPopUpModel{
		SpellNameInput:         spellNameInput,
		SpellCommandInput:      spellCmdInput,
		TotalInputField:        2,
		CurrentFocusInputIndex: 0,
		Error:                  nil,
		OnInputFuncTrigger:     learn.SaveSpell,
	}

	m.SpellPopUpModel = popUpModel
}
