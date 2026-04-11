package ui

import (
	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/textinput"
	"github.com/charmbracelet/x/ansi"
	"github.com/gohyuhan/rift/api/learn"
	apiUtils "github.com/gohyuhan/rift/api/utils"
	"github.com/gohyuhan/rift/i18n"
	"github.com/gohyuhan/rift/style"
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

func initCastLocationOptionPopUpModel(m *SpellbookInteractiveModel, spellName string) {
	castLocationOptionListArray := []list.Item{
		castLocationOptionItem{
			Title:       i18n.LANGUAGEMAPPING.CastLocationOptionCurrent,
			Description: i18n.LANGUAGEMAPPING.CastLocationOptionCurrentDescription,
			OptionType:  CastCWD,
		},
		castLocationOptionItem{
			Title:       i18n.LANGUAGEMAPPING.CastLocationOptionWaypoint,
			Description: i18n.LANGUAGEMAPPING.CastLocationOptionWaypointDescription,
			OptionType:  CastWaypoint,
		},
	}

	titleWidthLimit := m.Width - ListItemOrTitleWidthPad - ListTitleHorizontalPadding

	castLocationOptionList := list.New(castLocationOptionListArray, castLocationOptionDelegate{}, m.Width, m.Height)
	castLocationOptionList.SetShowPagination(false)
	castLocationOptionList.SetShowStatusBar(false)
	castLocationOptionList.SetFilteringEnabled(false)
	castLocationOptionList.SetShowFilter(false)
	castLocationOptionList.SetShowHelp(false)

	// truncate the title to prevent overflow when the terminal is narrow
	castLocationOptionList.Title = ansi.Truncate(i18n.LANGUAGEMAPPING.CastLocationOptionTitle, titleWidthLimit, "...")
	castLocationOptionList.Styles.Title = style.NewStyle.Bold(true)
	castLocationOptionList.Styles.PaginationStyle = style.NewStyle
	castLocationOptionList.Styles.TitleBar = style.NewStyle

	popUpModel := &CastLocationOptionPopUpModel{
		CastLocationOptionList: castLocationOptionList,
		SelectedSpellName:      spellName,
	}

	m.SpellPopUpModel = popUpModel
}

func initCastWaypointLocationOptionPopUpModel(m *SpellbookInteractiveModel, spellName string) {
	titleWidthLimit := m.Width - ListItemOrTitleWidthPad - ListTitleHorizontalPadding

	allWaypointsInfo, _ := apiUtils.GetAllWaypointsInfo(m.BboltReadDb)
	castWaypointLocationOptionListArray := []list.Item{}
	for _, waypoint := range allWaypointsInfo {
		castWaypointLocationOptionListArray = append(castWaypointLocationOptionListArray, castWaypointLocationOptionItem(waypoint))
	}

	castWaypointLocationOptionList := list.New(castWaypointLocationOptionListArray, castWaypointLocationDelegate{}, m.Width, m.Height)
	castWaypointLocationOptionList.SetShowPagination(false)
	castWaypointLocationOptionList.SetShowStatusBar(false)
	castWaypointLocationOptionList.SetFilteringEnabled(false)
	castWaypointLocationOptionList.SetShowFilter(false)
	castWaypointLocationOptionList.SetShowHelp(false)

	// truncate the title to prevent overflow when the terminal is narrow
	castWaypointLocationOptionList.Title = ansi.Truncate(i18n.LANGUAGEMAPPING.CastWaypointLocationOptionTitle, titleWidthLimit, "...")
	castWaypointLocationOptionList.Styles.Title = style.NewStyle.Bold(true)
	castWaypointLocationOptionList.Styles.PaginationStyle = style.NewStyle
	castWaypointLocationOptionList.Styles.TitleBar = style.NewStyle

	popUpModel := &CastWaypointLocationOptionPopUpModel{
		CastWaypointLocationOptionList: castWaypointLocationOptionList,
		SelectedSpellName:              spellName,
	}

	m.SpellPopUpModel = popUpModel
}
