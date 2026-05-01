package ui

import (
	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/gohyuhan/rift/api/inscribe"
	"github.com/gohyuhan/rift/api/scroll/features"
	apiUtils "github.com/gohyuhan/rift/api/utils"
	"github.com/gohyuhan/rift/i18n"
	"github.com/gohyuhan/rift/style"
)

// ----------------------------------
//
//	initialises the InscribePopUpModel with two text inputs (ritual name and ritual
//	description), configures placeholders from i18n, focuses the ritual name input,
//	and stores the populated model on the parent ScrollInteractiveModel;
//	OnInputFuncTrigger is set to learn.SaveRitual so submissions are persisted
//	to the database
//
// ----------------------------------
func initInscribePopUpModel(m *ScrollInteractiveModel, ritualName string, edit bool) tea.Cmd {
	totalInputField := 4 // default to 4 input fields (ritual name, description, commands, and submit button)
	ritualNameInput := textinput.New()
	ritualNameInput.SetValue("")
	ritualNameInput.Placeholder = i18n.LANGUAGEMAPPING.RitualNameInputPlaceHolder
	ritualNameInput.SetVirtualCursor(true)

	ritualDescInput := textarea.New()
	ritualDescInput.SetValue("")
	ritualDescInput.Placeholder = i18n.LANGUAGEMAPPING.RitualDescriptionInputPlaceHolder
	ritualDescInput.SetVirtualCursor(true)
	ritualDescInput.SetHeight(3)

	ritualCmdsInput := textarea.New()
	ritualCmdsInput.SetValue("")
	ritualCmdsInput.Placeholder = i18n.LANGUAGEMAPPING.RitualCommandsInputPlaceHolder
	ritualCmdsInput.SetVirtualCursor(true)
	ritualCmdsInput.SetHeight(7)

	var focusCmd tea.Cmd
	if edit {
		ritualInfoForEdit, viewErr := features.RetrieveRitualInfoDetailForEdit(ritualName)
		if viewErr == nil {
			ritualNameInput.SetValue(ritualInfoForEdit.RitualName)
			ritualDescInput.SetValue(ritualInfoForEdit.RitualDesc)
			ritualCmdsInput.SetValue(apiUtils.ParseRitualCommandsToString(ritualInfoForEdit.RitualCmds))
		}
		focusCmd = ritualDescInput.Focus()
		totalInputField = 3 // if in edit mode, only 3 input fields are needed since the ritual name cannot be changed
	} else {
		focusCmd = ritualNameInput.Focus()
	}

	popUpModel := &InscribePopUpModel{
		RitualName:             ritualName,
		RitualNameInput:        ritualNameInput,
		RitualDescriptionInput: ritualDescInput,
		RitualCommandsInput:    ritualCmdsInput,
		EditMode:               edit,
		TotalInputField:        totalInputField,
		CurrentFocusInputIndex: 0,
		Error:                  nil,
		OnInputFuncTrigger:     inscribe.SaveRitual,
	}

	m.RitualPopUpModel = popUpModel
	return focusCmd
}

// ----------------------------------
//
//	Builds the InvokeLocationOptionPopUpModel with two options — invoke at the
//	current working directory (InvokeCWD) or invoke at a waypoint (InvokeWaypoint).
//	The list title is truncated to prevent overflow on narrow terminals.
//	Stores the resulting model and the selected ritual name on the parent
//	ScrollInteractiveModel.
//
// ----------------------------------
func initInvokeLocationOptionPopUpModel(m *ScrollInteractiveModel, ritualName string) {
	invokeLocationOptionListArray := []list.Item{
		invokeLocationOptionItem{
			Title:       i18n.LANGUAGEMAPPING.InvokeLocationOptionCurrent,
			Description: i18n.LANGUAGEMAPPING.InvokeLocationOptionCurrentDescription,
			OptionType:  InvokeCWD,
		},
		invokeLocationOptionItem{
			Title:       i18n.LANGUAGEMAPPING.InvokeLocationOptionWaypoint,
			Description: i18n.LANGUAGEMAPPING.InvokeLocationOptionWaypointDescription,
			OptionType:  InvokeWaypoint,
		},
	}

	titleWidthLimit := m.Width - ListItemOrTitleWidthPad - ListTitleHorizontalPadding

	invokeLocationOptionList := list.New(invokeLocationOptionListArray, invokeLocationOptionDelegate{}, m.Width, m.Height)
	invokeLocationOptionList.SetShowPagination(false)
	invokeLocationOptionList.SetShowStatusBar(false)
	invokeLocationOptionList.SetFilteringEnabled(false)
	invokeLocationOptionList.SetShowFilter(false)
	invokeLocationOptionList.SetShowHelp(false)

	// truncate the title to prevent overflow when the terminal is narrow
	invokeLocationOptionList.Title = ansi.Truncate(i18n.LANGUAGEMAPPING.InvokeLocationOptionTitle, titleWidthLimit, "...")
	invokeLocationOptionList.Styles.Title = style.NewStyle.Bold(true)
	invokeLocationOptionList.Styles.PaginationStyle = style.NewStyle
	invokeLocationOptionList.Styles.TitleBar = style.NewStyle

	popUpModel := &InvokeLocationOptionPopUpModel{
		InvokeLocationOptionList: invokeLocationOptionList,
		SelectedRitualName:       ritualName,
	}

	m.RitualPopUpModel = popUpModel
}

// ----------------------------------
//
//	Builds the InvokeWaypointLocationOptionPopUpModel by loading all waypoints
//	from the DB and presenting them as selectable list items. The list title is
//	truncated to prevent overflow on narrow terminals. Stores the resulting model
//	and the selected ritual name on the parent ScrollInteractiveModel.
//
// ----------------------------------
func initInvokeWaypointLocationOptionPopUpModel(m *ScrollInteractiveModel, ritualName string) {
	titleWidthLimit := m.Width - ListItemOrTitleWidthPad - ListTitleHorizontalPadding

	allWaypointsInfo, _ := apiUtils.GetAllWaypointsInfo()
	invokeWaypointLocationOptionListArray := []list.Item{}
	for _, waypoint := range allWaypointsInfo {
		invokeWaypointLocationOptionListArray = append(invokeWaypointLocationOptionListArray, invokeWaypointLocationOptionItem(waypoint))
	}

	invokeWaypointLocationOptionList := list.New(invokeWaypointLocationOptionListArray, invokeWaypointLocationDelegate{}, m.Width, m.Height)
	invokeWaypointLocationOptionList.SetShowPagination(false)
	invokeWaypointLocationOptionList.SetShowStatusBar(false)
	invokeWaypointLocationOptionList.SetFilteringEnabled(false)
	invokeWaypointLocationOptionList.SetShowFilter(false)
	invokeWaypointLocationOptionList.SetShowHelp(false)

	// truncate the title to prevent overflow when the terminal is narrow
	invokeWaypointLocationOptionList.Title = ansi.Truncate(i18n.LANGUAGEMAPPING.InvokeWaypointLocationOptionTitle, titleWidthLimit, "...")
	invokeWaypointLocationOptionList.Styles.Title = style.NewStyle.Bold(true)
	invokeWaypointLocationOptionList.Styles.PaginationStyle = style.NewStyle
	invokeWaypointLocationOptionList.Styles.TitleBar = style.NewStyle

	popUpModel := &InvokeWaypointLocationOptionPopUpModel{
		InvokeWaypointLocationOptionList: invokeWaypointLocationOptionList,
		SelectedRitualName:               ritualName,
	}

	m.RitualPopUpModel = popUpModel
}
