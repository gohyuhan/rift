package ui

import (
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"github.com/gohyuhan/rift/api/waypoint/features"
	"github.com/gohyuhan/rift/i18n"
)

// ----------------------------------
//
//	initialises the rebind popup state on the model: creates a focused text
//	input pre-configured with the rebind placeholder, then stores a
//	RebindPopUpModel (containing the input, the target waypoint name, and the
//	RebindWaypoint trigger function) on m.WaypointPopUpModel ready for render
//
// ----------------------------------
func initRebindPopUp(m *WaypointInteractiveModel, waypointName string) tea.Cmd {
	rebindPathInput := textinput.New()
	rebindPathInput.SetValue("")
	rebindPathInput.Placeholder = i18n.LANGUAGEMAPPING.RebindPathInputPlaceHolder
	rebindPathInput.SetVirtualCursor(true)
	focusCmd := rebindPathInput.Focus()

	popUpModel := &RebindPopUpModel{
		RebindPathInput:    rebindPathInput,
		WaypointName:       waypointName,
		Error:              nil,
		OnInputFuncTrigger: features.RebindWaypoint,
	}

	m.WaypointPopUpModel = popUpModel
	return focusCmd
}

// ----------------------------------
//
//	initialises the reforge popup state on the model: creates a focused text
//	input pre-configured with the reforge placeholder, then stores a
//	ReforgePopUpModel (containing the input, the target waypoint name, and the
//	ReforgeWaypoint trigger function) on m.WaypointPopUpModel ready for render
//
// ----------------------------------
func initReforgePopUp(m *WaypointInteractiveModel, waypointName string) tea.Cmd {
	reforgeWaypointNameInput := textinput.New()
	reforgeWaypointNameInput.SetValue("")
	reforgeWaypointNameInput.Placeholder = i18n.LANGUAGEMAPPING.ReforgeWaypointNameInputPlaceHolder
	reforgeWaypointNameInput.SetVirtualCursor(true)
	focusCmd := reforgeWaypointNameInput.Focus()

	popUpModel := &ReforgePopUpModel{
		ReforgeWaypointNameInput: reforgeWaypointNameInput,
		WaypointName:             waypointName,
		Error:                    nil,
		OnInputFuncTrigger:       features.ReforgeWaypoint,
	}

	m.WaypointPopUpModel = popUpModel
	return focusCmd
}
