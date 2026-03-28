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
	RebindPathInput := textinput.New()
	RebindPathInput.SetValue("")
	RebindPathInput.Placeholder = i18n.LANGUAGEMAPPING.RebindPathInputPlaceHolder
	RebindPathInput.Focus()
	RebindPathInput.SetVirtualCursor(true)

	popUpModel := &RebindPopUpModel{
		RebindPathInput:    RebindPathInput,
		WaypointName:       waypointName,
		Error:              nil,
		OnInputFuncTrigger: features.RebindWaypoint,
	}

	m.WaypointPopUpModel = popUpModel
}
