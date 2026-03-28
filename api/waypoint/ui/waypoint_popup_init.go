package ui

import (
	"charm.land/bubbles/v2/textinput"
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
func initRebindPopUp(m *WaypointInteractiveModel, waypointName string) {
	rebindPathInput := textinput.New()
	rebindPathInput.SetValue("")
	rebindPathInput.Placeholder = i18n.LANGUAGEMAPPING.RebindPathInputPlaceHolder
	rebindPathInput.Focus()
	rebindPathInput.SetVirtualCursor(true)

	popUpModel := &RebindPopUpModel{
		RebindPathInput:    rebindPathInput,
		WaypointName:       waypointName,
		Error:              nil,
		OnInputFuncTrigger: features.RebindWaypoint,
	}

	m.WaypointPopUpModel = popUpModel
}

// ----------------------------------
//
//	initialises the reforge popup state on the model: creates a focused text
//	input pre-configured with the reforge placeholder, then stores a
//	ReforgePopUpModel (containing the input, the target waypoint name, and the
//	ReforgeWaypoint trigger function) on m.WaypointPopUpModel ready for render
//
// ----------------------------------
func initReforgePopUp(m *WaypointInteractiveModel, waypointName string) {
	reforgeWaypointNameInput := textinput.New()
	reforgeWaypointNameInput.SetValue("")
	reforgeWaypointNameInput.Placeholder = i18n.LANGUAGEMAPPING.ReforgeWaypointNameInputPlaceHolder
	reforgeWaypointNameInput.Focus()
	reforgeWaypointNameInput.SetVirtualCursor(true)

	popUpModel := &ReforgePopUpModel{
		ReforgeWaypointNameInput: reforgeWaypointNameInput,
		WaypointName:             waypointName,
		Error:                    nil,
		OnInputFuncTrigger:       features.ReforgeWaypoint,
	}

	m.WaypointPopUpModel = popUpModel
}
