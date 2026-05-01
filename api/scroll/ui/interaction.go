package ui

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/atotto/clipboard"
	"github.com/gohyuhan/rift/api/ritual"
	apiUtils "github.com/gohyuhan/rift/api/utils"
	"github.com/gohyuhan/rift/utils"
)

// ----------------------------------
//
//	processes key events when the model is not in typing mode;
//	j/↓ and k/↑ move the cursor, or scroll the help viewport when it is open;
//	n/N opens the inscribe popup when no popup is active; ? opens the help popup;
//	enter confirms the selected ritual; backspace forgets the selected ritual and
//	rebuilds the list
//
// ----------------------------------
func handleNonTypingInteraction(m *ScrollInteractiveModel, msg tea.KeyPressMsg) (*ScrollInteractiveModel, tea.Cmd) {
	switch msg.String() {
	case "j", "down":
		if m.ShowPopUp.Load() {
			var cmd tea.Cmd
			switch m.PopUpType {
			case HelpPopUp:
				m.RitualHelpViewport, cmd = m.RitualHelpViewport.Update(msg)
			case InvokeLocationOptionPopUp:
				popUp, ok := m.RitualPopUpModel.(*InvokeLocationOptionPopUpModel)
				if ok {
					popUp.InvokeLocationOptionList, cmd = popUp.InvokeLocationOptionList.Update(msg)
				}
			case InvokeWaypointLocationOptionPopUp:
				popUp, ok := m.RitualPopUpModel.(*InvokeWaypointLocationOptionPopUpModel)
				if ok {
					popUp.InvokeWaypointLocationOptionList, cmd = popUp.InvokeWaypointLocationOptionList.Update(msg)
				}
			}
			return m, cmd
		}
		m.RitualInfoList.CursorDown()
		m.RitualInfoListCursorPosition = m.RitualInfoList.Index()
	case "k", "up":
		if m.ShowPopUp.Load() {
			var cmd tea.Cmd
			switch m.PopUpType {
			case HelpPopUp:
				m.RitualHelpViewport, cmd = m.RitualHelpViewport.Update(msg)
			case InvokeLocationOptionPopUp:
				popUp, ok := m.RitualPopUpModel.(*InvokeLocationOptionPopUpModel)
				if ok {
					popUp.InvokeLocationOptionList, cmd = popUp.InvokeLocationOptionList.Update(msg)
				}
			case InvokeWaypointLocationOptionPopUp:
				popUp, ok := m.RitualPopUpModel.(*InvokeWaypointLocationOptionPopUpModel)
				if ok {
					popUp.InvokeWaypointLocationOptionList, cmd = popUp.InvokeWaypointLocationOptionList.Update(msg)
				}
			}
			return m, cmd
		}
		m.RitualInfoList.CursorUp()
		m.RitualInfoListCursorPosition = m.RitualInfoList.Index()
	case "n", "N":
		if !m.ShowPopUp.Load() {
			m.PopUpType = InscribePopUp
			m.ShowPopUp.Store(true)
			m.IsTypingMode.Store(true)
			return m, initInscribePopUpModel(m, "", false)
		}
	case "e", "E":
		if !m.ShowPopUp.Load() {
			currentSelectedRitual := m.RitualInfoList.SelectedItem()
			if currentSelectedRitual == nil {
				return m, nil
			}
			parsedRitualInfo := currentSelectedRitual.(ritualInfoItem)
			m.PopUpType = InscribePopUp
			m.ShowPopUp.Store(true)
			m.IsTypingMode.Store(true)
			return m, initInscribePopUpModel(m, parsedRitualInfo.RitualName, true)
		}
	case "?":
		m.ShowPopUp.Store(true)
		m.PopUpType = HelpPopUp
	case "enter":
		if m.ShowPopUp.Load() {
			switch m.PopUpType {
			case InvokeLocationOptionPopUp:
				popUp, ok := m.RitualPopUpModel.(*InvokeLocationOptionPopUpModel)
				if ok {
					selectedOption := popUp.InvokeLocationOptionList.SelectedItem()
					parsedSelectedOption := selectedOption.(invokeLocationOptionItem)
					switch parsedSelectedOption.OptionType {
					case InvokeCWD:
						executionPath, executionPathErr := utils.GetCWD()
						if executionPathErr != nil {
							executionPath = ""
						}
						m.SelectedRitualName = popUp.SelectedRitualName
						m.RitualInvokePath = executionPath
						m.IsQuit = true
						return m, tea.Quit
					case InvokeWaypoint:
						m.PopUpType = InvokeWaypointLocationOptionPopUp
						m.ShowPopUp.Store(true)
						initInvokeWaypointLocationOptionPopUpModel(m, popUp.SelectedRitualName)
						return m, nil
					}
				}
			case InvokeWaypointLocationOptionPopUp:
				popUp, ok := m.RitualPopUpModel.(*InvokeWaypointLocationOptionPopUpModel)
				if ok {
					selectedOption := popUp.InvokeWaypointLocationOptionList.SelectedItem()
					parsedSelectedWaypointOption := selectedOption.(invokeWaypointLocationOptionItem)
					m.SelectedRitualName = popUp.SelectedRitualName
					m.RitualInvokePath = parsedSelectedWaypointOption.WaypointPath
					m.IsQuit = true
					return m, tea.Quit
				}
			}
		} else {
			if ritualItem, ok := m.RitualInfoList.SelectedItem().(ritualInfoItem); ok {
				m.PopUpType = InvokeLocationOptionPopUp
				m.ShowPopUp.Store(true)
				initInvokeLocationOptionPopUpModel(m, ritualItem.RitualName)
				return m, nil
			}
		}
	case "backspace":
		if i, ok := m.RitualInfoList.SelectedItem().(ritualInfoItem); ok {
			// perform a forget on the selected ritual
			forgetErr := ritual.ForgetRitual(i.RitualName, false)
			if forgetErr != nil {
				m.ErrMessage = forgetErr
				m.IsQuit = true
				return m, tea.Quit
			}
			initRitualInfoListModel(m)
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
func handleTypingInteraction(m *ScrollInteractiveModel, msg tea.KeyPressMsg) (*ScrollInteractiveModel, tea.Cmd) {
	var cmd tea.Cmd
	switch msg.String() {
	case "ctrl+y":
		switch m.PopUpType {
		case InscribePopUp:
			popUp, ok := m.RitualPopUpModel.(*InscribePopUpModel)
			if ok {
				if popUp.EditMode {
					switch popUp.CurrentFocusInputIndex {
					case 0:
						clipboard.WriteAll(popUp.RitualDescriptionInput.Value())
					case 1:
						clipboard.WriteAll(popUp.RitualCommandsInput.Value())
					}
				} else {
					switch popUp.CurrentFocusInputIndex {
					case 0:
						clipboard.WriteAll(popUp.RitualNameInput.Value())
					case 1:
						clipboard.WriteAll(popUp.RitualDescriptionInput.Value())
					case 2:
						clipboard.WriteAll(popUp.RitualCommandsInput.Value())
					}
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
		case InscribePopUp:
			popUp, ok := m.RitualPopUpModel.(*InscribePopUpModel)
			if ok {
				if popUp.EditMode {
					switch popUp.CurrentFocusInputIndex {
					case 0:
						popUp.RitualDescriptionInput, cmd = popUp.RitualDescriptionInput.Update(msg)
					case 1:
						popUp.RitualCommandsInput, cmd = popUp.RitualCommandsInput.Update(msg)
						validateAndUpdateRitualPopUpState(popUp)
					}
				} else {
					switch popUp.CurrentFocusInputIndex {
					case 0:
						popUp.RitualNameInput, cmd = popUp.RitualNameInput.Update(msg)
					case 1:
						popUp.RitualDescriptionInput, cmd = popUp.RitualDescriptionInput.Update(msg)
					case 2:
						popUp.RitualCommandsInput, cmd = popUp.RitualCommandsInput.Update(msg)
						validateAndUpdateRitualPopUpState(popUp)
					}
				}
			}
		}
		return m, cmd

	case "enter":
		closePopUp := false

		switch m.PopUpType {
		case InscribePopUp:
			popUp, ok := m.RitualPopUpModel.(*InscribePopUpModel)
			if ok {
				if popUp.RitualCommandsInput.Focused() || popUp.RitualDescriptionInput.Focused() || popUp.RitualNameInput.Focused() || popUp.InscribeDisable.Load() {
					return m, nil
				}
				var ritualName string
				if popUp.EditMode {
					ritualName = popUp.RitualName
				} else {
					ritualName = strings.TrimSpace(popUp.RitualNameInput.Value())
				}
				ritualDesc := strings.TrimSpace(popUp.RitualDescriptionInput.Value())
				ritualCmds := strings.TrimSpace(popUp.RitualCommandsInput.Value())
				inscribeErr := popUp.OnInputFuncTrigger(ritualName, ritualDesc, ritualCmds, popUp.EditMode)
				if inscribeErr != nil {
					// surface the error in the popup rather than quitting
					popUp.Error = inscribeErr
					return m, nil
				}
				closePopUp = true
				initRitualInfoListModel(m)
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
		case InscribePopUp:
			popUp, ok := m.RitualPopUpModel.(*InscribePopUpModel)
			if ok {
				popUp.RitualNameInput.Blur()
				popUp.RitualDescriptionInput.Blur()
				popUp.RitualCommandsInput.Blur()
				if popUp.EditMode {
					switch popUp.CurrentFocusInputIndex {
					case 0:
						popUp.CurrentFocusInputIndex = 1
						cmd = popUp.RitualCommandsInput.Focus()
					case 1:
						popUp.CurrentFocusInputIndex = 2
					case 2:
						// does nothing
					}
				} else {
					switch popUp.CurrentFocusInputIndex {
					case 0:
						popUp.CurrentFocusInputIndex = 1
						cmd = popUp.RitualDescriptionInput.Focus()
					case 1:
						popUp.CurrentFocusInputIndex = 2
						cmd = popUp.RitualCommandsInput.Focus()
					case 2:
						popUp.CurrentFocusInputIndex = 3
					case 3:
						// does nothing
					}
				}
				return m, cmd
			}
		}

	case "shift+tab":
		switch m.PopUpType {
		case InscribePopUp:
			popUp, ok := m.RitualPopUpModel.(*InscribePopUpModel)
			if ok {
				popUp.RitualNameInput.Blur()
				popUp.RitualDescriptionInput.Blur()
				popUp.RitualCommandsInput.Blur()
				if popUp.EditMode {
					switch popUp.CurrentFocusInputIndex {
					case 0:
						// does nothing
					case 1:
						popUp.CurrentFocusInputIndex = 0
						cmd = popUp.RitualDescriptionInput.Focus()
					case 2:
						popUp.CurrentFocusInputIndex = 1
						cmd = popUp.RitualCommandsInput.Focus()
					}
				} else {
					switch popUp.CurrentFocusInputIndex {
					case 0:
						// does nothing
					case 1:
						popUp.CurrentFocusInputIndex = 0
						cmd = popUp.RitualNameInput.Focus()
					case 2:
						popUp.CurrentFocusInputIndex = 1
						cmd = popUp.RitualDescriptionInput.Focus()
					case 3:
						popUp.CurrentFocusInputIndex = 2
						cmd = popUp.RitualCommandsInput.Focus()
					}
				}
				return m, cmd
			}
		}
	}

	// forward all other key events to the active popup's input component
	switch m.PopUpType {
	case InscribePopUp:
		popUp, ok := m.RitualPopUpModel.(*InscribePopUpModel)
		if ok {
			if popUp.EditMode {
				switch popUp.CurrentFocusInputIndex {
				case 0:
					popUp.RitualDescriptionInput, cmd = popUp.RitualDescriptionInput.Update(msg)
				case 1:
					popUp.RitualCommandsInput, cmd = popUp.RitualCommandsInput.Update(msg)
					validateAndUpdateRitualPopUpState(popUp)
				}
			} else {
				switch popUp.CurrentFocusInputIndex {
				case 0:
					popUp.RitualNameInput, cmd = popUp.RitualNameInput.Update(msg)
				case 1:
					popUp.RitualDescriptionInput, cmd = popUp.RitualDescriptionInput.Update(msg)
				case 2:
					popUp.RitualCommandsInput, cmd = popUp.RitualCommandsInput.Update(msg)
					validateAndUpdateRitualPopUpState(popUp)
				}
			}
			return m, cmd
		}
	}
	return m, nil
}

func validateAndUpdateRitualPopUpState(popUp *InscribePopUpModel) {
	normCmds, err := apiUtils.NormalizeAndCheckRitualCommandsAreValid(popUp.RitualCommandsInput.Value())
	if err != nil {
		popUp.Error = err
		popUp.InscribeDisable.Store(true)
	} else {
		popUp.Error = nil
		popUp.InscribeDisable.Store(len(normCmds) == 0)
	}
}
