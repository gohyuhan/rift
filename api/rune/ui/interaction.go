package ui

import (
	tea "charm.land/bubbletea/v2"
	"github.com/atotto/clipboard"
	apiUtils "github.com/gohyuhan/rift/api/utils"
)

// ----------------------------------
//
//	handles key events when the rune UI is in navigation mode (not typing).
//	j/k/up/down scroll the ChooseRuneEngraveOptionPopUp list.
//	enter on the option list either transitions to the EngraveRuneCommandsPopUp
//	(for enter/leave engrave options) or immediately removes the rune slot
//	(for remove enter/leave options) and refreshes the option list.
//
// ----------------------------------
func handleNonTypingInteraction(m *RuneInteractiveModel, msg tea.KeyPressMsg) (*RuneInteractiveModel, tea.Cmd) {
	switch msg.String() {
	case "j", "down", "k", "up":
		if m.ShowPopUp.Load() {
			var cmd tea.Cmd
			switch m.PopUpType {
			case ChooseRuneEngraveOptionPopUp:
				popUp, ok := m.RunePopUpModel.(*ChooseRuneEngraveOptionPopUpModel)
				if ok {
					popUp.EngraveOptionList, cmd = popUp.EngraveOptionList.Update(msg)
				}
			}
			return m, cmd
		}
	case "enter":
		if m.ShowPopUp.Load() {
			var cmd tea.Cmd
			switch m.PopUpType {
			case ChooseRuneEngraveOptionPopUp:
				popUp, ok := m.RunePopUpModel.(*ChooseRuneEngraveOptionPopUpModel)
				if ok {
					selectedOption := popUp.EngraveOptionList.SelectedItem()
					parsedSelectedOption := selectedOption.(runeEngraveOptionInfoItem)
					switch parsedSelectedOption.runeEngraveOptionType {
					case EngraveRuneEnterType:
						m.PopUpType = EngraveRuneCommandsPopUp
						m.IsTypingMode.Store(true)
						initEngraveRuneCommandsPopUpModel(m, EngraveRuneEnterType)
					case EngraveRuneLeaveType:
						m.PopUpType = EngraveRuneCommandsPopUp
						m.IsTypingMode.Store(true)
						initEngraveRuneCommandsPopUpModel(m, EngraveRuneLeaveType)
					case RemoveRuneEnterType:
						apiUtils.RemoveEnterRuneCmds(m.ChosenWaypointPath)
						cmd = initChooseRuneEngraveOptionPopUpModel(m)
					case RemoveRuneLeaveType:
						apiUtils.RemoveLeaveRuneCmds(m.ChosenWaypointPath)
						cmd = initChooseRuneEngraveOptionPopUpModel(m)
					}
				}
			}
			return m, cmd
		}
	}
	return m, nil
}

// ----------------------------------
//
//	handles key events when the rune UI is in typing mode (textarea focused or
//	button selected in EngraveRuneCommandsPopUp).
//	ctrl+y copies the textarea content to the clipboard;
//	ctrl+p pastes clipboard content into the textarea.
//	tab moves focus from the textarea to the Engrave button;
//	shift+tab moves focus back to the textarea.
//	enter on the Engrave button (when enabled) validates and persists the rune
//	commands, sets RuneEngraved, and quits the program.
//	all other keys are forwarded to the textarea when it has focus.
//
// ----------------------------------
func handleTypingInteraction(m *RuneInteractiveModel, msg tea.KeyPressMsg) (*RuneInteractiveModel, tea.Cmd) {
	var cmd tea.Cmd
	switch msg.String() {
	case "ctrl+y":
		switch m.PopUpType {
		case EngraveRuneCommandsPopUp:
			popUp, ok := m.RunePopUpModel.(*EngraveRuneCommandsPopUpModel)
			if ok && popUp.TextAreaFocused.Load() {
				clipboard.WriteAll(popUp.RuneCommandsTextArea.Value())
			}
		}
		return m, cmd

	case "ctrl+p":
		content, err := clipboard.ReadAll()
		if err != nil {
			return m, nil
		}
		msg := tea.PasteMsg{
			Content: content,
		}
		switch m.PopUpType {
		case EngraveRuneCommandsPopUp:
			popUp, ok := m.RunePopUpModel.(*EngraveRuneCommandsPopUpModel)
			if ok && popUp.TextAreaFocused.Load() {
				popUp.RuneCommandsTextArea, cmd = popUp.RuneCommandsTextArea.Update(msg)
				validateAndUpdateRunePopUpState(popUp)
			}
		}
		return m, cmd

	case "enter":
		popUp, ok := m.RunePopUpModel.(*EngraveRuneCommandsPopUpModel)
		if ok && !popUp.TextAreaFocused.Load() && !popUp.EngraveDisable.Load() {
			normRuneCmd, runeErr := apiUtils.NormalizeAndCheckRuneCommandsAreValid(popUp.RuneCommandsTextArea.Value())
			if runeErr != nil {
				m.ErrMessage = runeErr
			} else {
				var engraveErr error
				switch popUp.RuneEngraveOptionType {
				case EngraveRuneEnterType:
					engraveErr = apiUtils.EngraveEnterRuneCmds(m.ChosenWaypointPath, normRuneCmd)
				case EngraveRuneLeaveType:
					engraveErr = apiUtils.EngraveLeaveRuneCmds(m.ChosenWaypointPath, normRuneCmd)
				}

				if engraveErr != nil {
					m.ErrMessage = engraveErr
				} else {
					m.ErrMessage = nil
					m.RuneEngraved.Store(true)
				}
			}
			m.IsQuit = true
			return m, tea.Quit
		}

	case "tab":
		switch m.PopUpType {
		case EngraveRuneCommandsPopUp:
			popUp, ok := m.RunePopUpModel.(*EngraveRuneCommandsPopUpModel)
			if ok {
				popUp.TextAreaFocused.Store(false)
				return m, cmd
			}
		}

	case "shift+tab":
		switch m.PopUpType {
		case EngraveRuneCommandsPopUp:
			popUp, ok := m.RunePopUpModel.(*EngraveRuneCommandsPopUpModel)
			if ok {
				popUp.TextAreaFocused.Store(true)
				return m, cmd
			}
		}
	}

	// forward all other key events to the active popup's input component
	switch m.PopUpType {
	case EngraveRuneCommandsPopUp:
		popUp, ok := m.RunePopUpModel.(*EngraveRuneCommandsPopUpModel)
		if ok {
			if popUp.TextAreaFocused.Load() {
				popUp.RuneCommandsTextArea, cmd = popUp.RuneCommandsTextArea.Update(msg)
				validateAndUpdateRunePopUpState(popUp)
			}
			return m, cmd
		}
	}
	return m, nil
}

func validateAndUpdateRunePopUpState(popUp *EngraveRuneCommandsPopUpModel) {
	normCmds, err := apiUtils.NormalizeAndCheckRuneCommandsAreValid(popUp.RuneCommandsTextArea.Value())
	if err != nil {
		popUp.Error = err
		popUp.EngraveDisable.Store(true)
	} else {
		popUp.Error = nil
		popUp.EngraveDisable.Store(len(normCmds) == 0)
	}
}
