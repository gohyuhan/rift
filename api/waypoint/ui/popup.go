package ui

import (
	"charm.land/lipgloss/v2"
	"github.com/gohyuhan/rift/i18n"
	"github.com/gohyuhan/rift/style"
)

// ----------------------------------
//
//	dispatches to the appropriate popup renderer based on the active PopUpType;
//	returns an empty string if the type is unrecognised or not yet implemented
//
// ----------------------------------
func renderPopUpComponent(m *WaypointInteractiveModel) string {
	switch m.PopUpType {
	case HelpPopUp:
		return renderHelpPopUp(m)
	case RebindPopUp:
		return renderRebindPopUp(m)
	case ReforgePopUp:
		return renderReforgePopUp(m)
	}
	return ""
}

// ----------------------------------
//
//	renders the help popup: sizes the viewport to 80 % width / 70 % height of
//	the terminal, then wraps it in a rounded border
//
// ----------------------------------
func renderHelpPopUp(m *WaypointInteractiveModel) string {
	maxWidth := int(float64(m.Width) * 0.8)
	maxHeight := int(float64(m.Height) * 0.7)
	m.WaypointHelpViewport.SetWidth(maxWidth - ListItemOrTitleWidthPad)
	m.WaypointHelpViewport.SetHeight(maxHeight - ListItemOrTitleWidthPad)
	return style.BorderStyle.
		MaxWidth(maxWidth).
		MaxHeight(maxHeight).
		Render(m.WaypointHelpViewport.View())
}

// ----------------------------------
//
//	renders the rebind path input popup: sizes the text input to 80 % width,
//	stacks the title and input vertically, appends a coloured error line when
//	the last submit attempt failed, then wraps the whole block in a rounded
//	border; returns an empty string if the popup model is the wrong type
//
// ----------------------------------
func renderRebindPopUp(m *WaypointInteractiveModel) string {
	popUp, ok := m.WaypointPopUpModel.(*RebindPopUpModel)
	if ok {
		maxWidth := int(float64(m.Width) * 0.8)
		maxHeight := int(float64(m.Height) * 0.7)
		popUp.RebindPathInput.SetWidth(maxWidth - TextInputWidthPad)
		strStyle := lipgloss.NewStyle().Width(maxWidth - TextInputWidthPad)
		var content string
		if popUp.Error != nil {
			errMessage := strStyle.Render(style.RenderStringWithColor(popUp.Error.Error(), style.ColorError, true))
			content = lipgloss.JoinVertical(
				lipgloss.Left,
				strStyle.Render(i18n.LANGUAGEMAPPING.WaypointRebindTitle),
				popUp.RebindPathInput.View(),
				errMessage,
			)
		} else {
			content = lipgloss.JoinVertical(
				lipgloss.Left,
				strStyle.Render(i18n.LANGUAGEMAPPING.WaypointRebindTitle),
				popUp.RebindPathInput.View(),
			)
		}
		return style.BorderStyle.
			MaxWidth(maxWidth).
			MaxHeight(maxHeight).
			Render(content)
	}
	return ""
}

// ----------------------------------
//
//	renders the reforge name input popup: sizes the text input to 80 % width,
//	stacks the title and input vertically, appends a coloured error line when
//	the last submit attempt failed, then wraps the whole block in a rounded
//	border; returns an empty string if the popup model is the wrong type
//
// ----------------------------------
func renderReforgePopUp(m *WaypointInteractiveModel) string {
	popUp, ok := m.WaypointPopUpModel.(*ReforgePopUpModel)
	if ok {
		maxWidth := int(float64(m.Width) * 0.8)
		maxHeight := int(float64(m.Height) * 0.7)
		popUp.ReforgeWaypointNameInput.SetWidth(maxWidth - TextInputWidthPad)
		strStyle := lipgloss.NewStyle().Width(maxWidth - TextInputWidthPad)
		var content string
		if popUp.Error != nil {
			errMessage := strStyle.Render(style.RenderStringWithColor(popUp.Error.Error(), style.ColorError, true))
			content = lipgloss.JoinVertical(
				lipgloss.Left,
				strStyle.Render(i18n.LANGUAGEMAPPING.WaypointReforgeTitle),
				popUp.ReforgeWaypointNameInput.View(),
				errMessage,
			)
		} else {
			content = lipgloss.JoinVertical(
				lipgloss.Left,
				strStyle.Render(i18n.LANGUAGEMAPPING.WaypointReforgeTitle),
				popUp.ReforgeWaypointNameInput.View(),
			)
		}
		return style.BorderStyle.
			MaxWidth(maxWidth).
			MaxHeight(maxHeight).
			Render(content)
	}
	return ""
}
