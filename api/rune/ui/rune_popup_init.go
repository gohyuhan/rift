package ui

import (
	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/textarea"
	tea "charm.land/bubbletea/v2"
	apiUtils "github.com/gohyuhan/rift/api/utils"
	"github.com/gohyuhan/rift/i18n"
	"github.com/gohyuhan/rift/style"
)

// ----------------------------------
//
//	Fetches the waypoint's current rune state and builds the
//	ChooseRuneEngraveOptionPopUp list. Always includes engrave enter and leave
//	options; conditionally appends remove enter and remove leave options when
//	the corresponding rune slot is already populated. Stores the resulting
//	ChooseRuneEngraveOptionPopUpModel on the parent RuneInteractiveModel and
//	returns nil (no startup command needed). Quits the program on DB error.
//
// ----------------------------------
func initChooseRuneEngraveOptionPopUpModel(m *RuneInteractiveModel) tea.Cmd {
	wp, err := apiUtils.RetrieveWaypointInfo(m.ChosenWaypointName)
	if err != nil {
		m.ErrMessage = err
		m.IsQuit = true
		return tea.Quit
	}
	m.ChosenWaypointPath = wp.WaypointPath

	m.ExistingEnterRune = wp.EnterRunes
	m.ExistingLeaveRune = wp.LeaveRunes

	engraveOptionListArray := []list.Item{
		runeEngraveOptionInfoItem{
			runeEngraveOptionName: i18n.LANGUAGEMAPPING.EngraveRuneEnterOptionName,
			runeEngraveOptionDesc: i18n.LANGUAGEMAPPING.EngraveRuneEnterOptionDesc,
			runeEngraveOptionType: EngraveRuneEnterType,
		},
		runeEngraveOptionInfoItem{
			runeEngraveOptionName: i18n.LANGUAGEMAPPING.EngraveRuneLeaveOptionName,
			runeEngraveOptionDesc: i18n.LANGUAGEMAPPING.EngraveRuneLeaveOptionDesc,
			runeEngraveOptionType: EngraveRuneLeaveType,
		},
	}

	if m.ExistingEnterRune != nil || len(m.ExistingEnterRune) > 0 {
		engraveOptionListArray = append(engraveOptionListArray, runeEngraveOptionInfoItem{
			runeEngraveOptionName: i18n.LANGUAGEMAPPING.RemoveRuneEnterOptionName,
			runeEngraveOptionDesc: i18n.LANGUAGEMAPPING.RemoveRuneEnterOptionDesc,
			runeEngraveOptionType: RemoveRuneEnterType,
		})
	}

	if m.ExistingLeaveRune != nil || len(m.ExistingLeaveRune) > 0 {
		engraveOptionListArray = append(engraveOptionListArray, runeEngraveOptionInfoItem{
			runeEngraveOptionName: i18n.LANGUAGEMAPPING.RemoveRuneLeaveOptionName,
			runeEngraveOptionDesc: i18n.LANGUAGEMAPPING.RemoveRuneLeaveOptionDesc,
			runeEngraveOptionType: RemoveRuneLeaveType,
		})
	}

	engraveOptionList := list.New(engraveOptionListArray, runeEngraveOptionInfoDelegate{}, m.Width, m.Height)
	engraveOptionList.SetShowTitle(false)
	engraveOptionList.SetShowPagination(false)
	engraveOptionList.SetShowStatusBar(false)
	engraveOptionList.SetFilteringEnabled(false)
	engraveOptionList.SetShowFilter(false)
	engraveOptionList.SetShowHelp(false)

	// truncate the title to prevent overflow when the terminal is narrow
	engraveOptionList.Styles.PaginationStyle = style.NewStyle

	popUpModel := &ChooseRuneEngraveOptionPopUpModel{
		EngraveOptionList: engraveOptionList,
	}

	m.RunePopUpModel = popUpModel
	return nil
}

// ----------------------------------
//
//	Builds the EngraveRuneCommandsPopUpModel for the given rune slot
//	(EngraveRuneEnterType or EngraveRuneLeaveType). Pre-populates the textarea
//	with any existing rune commands for that slot (formatted as a newline-
//	separated string), focuses the textarea, and stores the model on the parent
//	RuneInteractiveModel. The Engrave button starts disabled; it is enabled once
//	the textarea contains at least one valid non-cd command.
//
// ----------------------------------
func initEngraveRuneCommandsPopUpModel(m *RuneInteractiveModel, runeEngraveOptionType string) {
	var value string
	switch runeEngraveOptionType {
	case EngraveRuneEnterType:
		value = apiUtils.ParseRuneCommandsToString(m.ExistingEnterRune)
	case EngraveRuneLeaveType:
		value = apiUtils.ParseRuneCommandsToString(m.ExistingLeaveRune)
	}
	runeCommandsTextArea := textarea.New()
	runeCommandsTextArea.SetValue(value)
	runeCommandsTextArea.Placeholder = i18n.LANGUAGEMAPPING.RuneCommandsPlaceHolder
	runeCommandsTextArea.Focus()
	runeCommandsTextArea.SetVirtualCursor(true)

	popUpModel := &EngraveRuneCommandsPopUpModel{
		RuneCommandsTextArea:  runeCommandsTextArea,
		RuneEngraveOptionType: runeEngraveOptionType,
		Error:                 nil,
	}

	popUpModel.EngraveDisable.Store(true)
	popUpModel.TextAreaFocused.Store(true)

	m.RunePopUpModel = popUpModel
}
