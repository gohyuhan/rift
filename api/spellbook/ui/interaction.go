package ui

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/atotto/clipboard"
	"github.com/gohyuhan/rift/api/spell"
	"github.com/gohyuhan/rift/utils"
)

// ----------------------------------
//
//	processes key events when the model is not in typing mode;
//	j/↓ and k/↑ move the cursor, or scroll the help viewport when it is open;
//	n/N opens the learn popup when no popup is active; ? opens the help popup;
//	enter confirms the selected spell; backspace forgets the selected spell and
//	rebuilds the list
//
// ----------------------------------
func handleNonTypingInteraction(m *SpellbookInteractiveModel, msg tea.KeyPressMsg) (*SpellbookInteractiveModel, tea.Cmd) {
	switch msg.String() {
	case "j", "down":
		if m.ShowPopUp.Load() {
			var cmd tea.Cmd
			switch m.PopUpType {
			case HelpPopUp:
				m.SpellHelpViewport, cmd = m.SpellHelpViewport.Update(msg)
			case CastLocationOptionPopUp:
				popUp, ok := m.SpellPopUpModel.(*CastLocationOptionPopUpModel)
				if ok {
					popUp.CastLocationOptionList, cmd = popUp.CastLocationOptionList.Update(msg)
				}
			case CastWaypointLocationOptionPopUp:
				popUp, ok := m.SpellPopUpModel.(*CastWaypointLocationOptionPopUpModel)
				if ok {
					popUp.CastWaypointLocationOptionList, cmd = popUp.CastWaypointLocationOptionList.Update(msg)
				}
			}
			return m, cmd
		}
		m.SpellInfoList.CursorDown()
		m.SpellInfoListCursorPosition = m.SpellInfoList.Index()
	case "k", "up":
		if m.ShowPopUp.Load() {
			var cmd tea.Cmd
			switch m.PopUpType {
			case HelpPopUp:
				m.SpellHelpViewport, cmd = m.SpellHelpViewport.Update(msg)
			case CastLocationOptionPopUp:
				popUp, ok := m.SpellPopUpModel.(*CastLocationOptionPopUpModel)
				if ok {
					popUp.CastLocationOptionList, cmd = popUp.CastLocationOptionList.Update(msg)
				}
			case CastWaypointLocationOptionPopUp:
				popUp, ok := m.SpellPopUpModel.(*CastWaypointLocationOptionPopUpModel)
				if ok {
					popUp.CastWaypointLocationOptionList, cmd = popUp.CastWaypointLocationOptionList.Update(msg)
				}
			}
			return m, cmd
		}
		m.SpellInfoList.CursorUp()
		m.SpellInfoListCursorPosition = m.SpellInfoList.Index()
	case "n", "N":
		if !m.ShowPopUp.Load() {
			m.PopUpType = LearnPopUp
			m.ShowPopUp.Store(true)
			m.IsTypingMode.Store(true)
			initLearnPopUpModel(m)
		}
	case "?":
		m.ShowPopUp.Store(true)
		m.PopUpType = HelpPopUp
	case "enter":
		if m.ShowPopUp.Load() {
			switch m.PopUpType {
			case CastLocationOptionPopUp:
				popUp, ok := m.SpellPopUpModel.(*CastLocationOptionPopUpModel)
				if ok {
					selectedOption := popUp.CastLocationOptionList.SelectedItem()
					parsedSelectedOption := selectedOption.(castLocationOptionItem)
					switch parsedSelectedOption.OptionType {
					case CastCWD:
						executionPath, executionPathErr := utils.GetCWD()
						if executionPathErr != nil {
							executionPath = ""
						}
						m.SelectedSpellName = popUp.SelectedSpellName
						m.SpellCastPath = executionPath
						m.IsQuit = true
						return m, tea.Quit
					case CastWaypoint:
						// TODO: implement cast to waypoint functionality
						m.PopUpType = CastWaypointLocationOptionPopUp
						m.ShowPopUp.Store(true)
						initCastWaypointLocationOptionPopUpModel(m, popUp.SelectedSpellName)
						return m, nil
					}
				}
			case CastWaypointLocationOptionPopUp:
				popUp, ok := m.SpellPopUpModel.(*CastWaypointLocationOptionPopUpModel)
				if ok {
					selectedOption := popUp.CastWaypointLocationOptionList.SelectedItem()
					parsedSelectedWaypointOption := selectedOption.(castWaypointLocationOptionItem)
					m.SelectedSpellName = popUp.SelectedSpellName
					m.SpellCastPath = parsedSelectedWaypointOption.WaypointPath
					m.IsQuit = true
					return m, tea.Quit
				}
			}
		} else {
			if spellItem, ok := m.SpellInfoList.SelectedItem().(spellInfoItem); ok {
				m.PopUpType = CastLocationOptionPopUp
				m.ShowPopUp.Store(true)
				initCastLocationOptionPopUpModel(m, spellItem.SpellName)
				return m, nil
			}
		}
	case "backspace":
		if i, ok := m.SpellInfoList.SelectedItem().(spellInfoItem); ok {
			// perform a forget on the selected spell
			forgetErr := spell.ForgetSpell(m.BboltDb, i.SpellName, false)
			if forgetErr != nil {
				m.ErrMessage = forgetErr
				m.IsQuit = true
				return m, tea.Quit
			}
			initSpellInfoListModel(m)
			return m, nil
		}
	}
	return m, nil
}

// ----------------------------------
//
//	processes key events when the model is in typing mode (a text input popup
//	is active); ctrl+y copies the focused input's value to the clipboard;
//	ctrl+p pastes clipboard content into the focused input; on enter, submits
//	the active popup's inputs by calling the stored OnInputFuncTrigger and
//	closes the popup on success, or surfaces the error inline on failure;
//	tab advances focus to the next input field (no-op on the last field);
//	shift+tab moves focus to the previous input field (no-op on the first);
//	all other keys are forwarded to the focused input component
//
// ----------------------------------
func handleTypingInteraction(m *SpellbookInteractiveModel, msg tea.KeyPressMsg) (*SpellbookInteractiveModel, tea.Cmd) {
	var cmd tea.Cmd
	switch msg.String() {
	case "ctrl+y":
		switch m.PopUpType {
		case LearnPopUp:
			popUp, ok := m.SpellPopUpModel.(*LearnPopUpModel)
			if ok {
				switch popUp.CurrentFocusInputIndex {
				case 0:
					clipboard.WriteAll(popUp.SpellNameInput.Value())
				case 1:
					clipboard.WriteAll(popUp.SpellCommandInput.Value())
				}
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
		case LearnPopUp:
			popUp, ok := m.SpellPopUpModel.(*LearnPopUpModel)
			if ok {
				switch popUp.CurrentFocusInputIndex {
				case 0:
					popUp.SpellNameInput, cmd = popUp.SpellNameInput.Update(msg)
				case 1:
					popUp.SpellCommandInput, cmd = popUp.SpellCommandInput.Update(msg)
				}
			}
		}
		return m, cmd

	case "enter":
		closePopUp := false

		switch m.PopUpType {
		case LearnPopUp:
			popUp, ok := m.SpellPopUpModel.(*LearnPopUpModel)
			if ok {
				spellName := strings.TrimSpace(popUp.SpellNameInput.Value())
				spellCmd := strings.TrimSpace(popUp.SpellCommandInput.Value())
				_, learnErr := popUp.OnInputFuncTrigger(m.BboltDb, spellName, spellCmd)
				if learnErr != nil {
					// surface the error in the popup rather than quitting
					popUp.Error = learnErr
					return m, nil
				}
				closePopUp = true
				initSpellInfoListModel(m)
			}

			if closePopUp {
				m.ShowPopUp.Store(false)
				m.IsTypingMode.Store(false)
				m.PopUpType = NoPopUp
			}
			return m, nil
		}

	case "tab":
		switch m.PopUpType {
		case LearnPopUp:
			popUp, ok := m.SpellPopUpModel.(*LearnPopUpModel)
			if ok {
				switch popUp.CurrentFocusInputIndex {
				case 0:
					popUp.CurrentFocusInputIndex = 1
					popUp.SpellNameInput.Blur()
					popUp.SpellCommandInput.Focus()
				case 1:
					// does nothing
				}
				return m, cmd
			}
		}

	case "shift+tab":
		switch m.PopUpType {
		case LearnPopUp:
			popUp, ok := m.SpellPopUpModel.(*LearnPopUpModel)
			if ok {
				switch popUp.CurrentFocusInputIndex {
				case 0:
					// does nothing
				case 1:
					popUp.CurrentFocusInputIndex = 0
					popUp.SpellNameInput.Focus()
					popUp.SpellCommandInput.Blur()
				}
				return m, cmd
			}
		}
	}

	// forward all other key events to the active popup's input component
	switch m.PopUpType {
	case LearnPopUp:
		popUp, ok := m.SpellPopUpModel.(*LearnPopUpModel)
		if ok {
			switch popUp.CurrentFocusInputIndex {
			case 0:
				popUp.SpellNameInput, cmd = popUp.SpellNameInput.Update(msg)
			case 1:
				popUp.SpellCommandInput, cmd = popUp.SpellCommandInput.Update(msg)
			}
			return m, cmd
		}
	}
	return m, nil
}
