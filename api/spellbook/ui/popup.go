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
func renderPopUpComponent(m *SpellbookInteractiveModel) string {
	switch m.PopUpType {
	case HelpPopUp:
		return renderHelpPopUp(m)
	case LearnPopUp:
		return renderLearnPopUp(m)
	case CastLocationOptionPopUp:
		return renderCastLocationOptionPopUp(m)
	case CastWaypointLocationOptionPopUp:
		return renderCastWaypointLocationOptionPopUp(m)
	}
	return ""
}

// ----------------------------------
//
//	renders the help popup: sizes the viewport to 80 % width / 70 % height of
//	the terminal, then wraps it in a rounded border
//
// ----------------------------------
func renderHelpPopUp(m *SpellbookInteractiveModel) string {
	maxWidth := int(float64(m.Width) * 0.8)
	maxHeight := int(float64(m.Height) * 0.7)
	m.SpellHelpViewport.SetWidth(maxWidth - ListItemOrTitleWidthPad)
	m.SpellHelpViewport.SetHeight(maxHeight - ListItemOrTitleWidthPad)
	return style.BorderStyle.
		MaxWidth(maxWidth).
		MaxHeight(maxHeight).
		Render(m.SpellHelpViewport.View())
}

// ----------------------------------
//
//	renders the learn new spell popup: casts SpellPopUpModel to *LearnPopUpModel
//	and sizes both text inputs to 80 % of the terminal width; joins the spell
//	name title, spell name input, command title, and command input vertically —
//	appending an inline error message below the inputs when one is present;
//	wraps the result in a rounded border sized to 80 % width / 70 % height;
//	returns an empty string if the model type assertion fails
//
// ----------------------------------
func renderLearnPopUp(m *SpellbookInteractiveModel) string {
	popUp, ok := m.SpellPopUpModel.(*LearnPopUpModel)
	if ok {
		maxWidth := int(float64(m.Width) * 0.8)
		maxHeight := int(float64(m.Height) * 0.7)
		popUp.SpellNameInput.SetWidth(maxWidth - TextInputWidthPad)
		popUp.SpellCommandInput.SetWidth(maxWidth - TextInputWidthPad)
		strStyle := lipgloss.NewStyle().Width(maxWidth - TextInputWidthPad)
		var content string
		if popUp.Error != nil {
			errMessage := strStyle.Render(style.RenderStringWithColor(popUp.Error.Error(), style.ColorError, true))
			content = lipgloss.JoinVertical(
				lipgloss.Left,
				strStyle.Render(i18n.LANGUAGEMAPPING.SpellNameInputTitle),
				popUp.SpellNameInput.View(),
				strStyle.Render(i18n.LANGUAGEMAPPING.SpellCommandInputTitle),
				popUp.SpellCommandInput.View(),
				errMessage,
			)
		} else {
			content = lipgloss.JoinVertical(
				lipgloss.Left,
				strStyle.Render(i18n.LANGUAGEMAPPING.SpellNameInputTitle),
				popUp.SpellNameInput.View(),
				strStyle.Render(i18n.LANGUAGEMAPPING.SpellCommandInputTitle),
				popUp.SpellCommandInput.View(),
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
//	renders the cast-location option popup: sizes the list to 80 % of the
//	terminal width and a fixed 12-row height, then wraps it in a rounded border;
//	returns an empty string if the popup model is the wrong type
//
// ----------------------------------
func renderCastLocationOptionPopUp(m *SpellbookInteractiveModel) string {
	popUp, ok := m.SpellPopUpModel.(*CastLocationOptionPopUpModel)
	if ok {
		maxWidth := int(float64(m.Width) * 0.8)
		maxHeight := 12
		popUp.CastLocationOptionList.SetWidth(maxWidth - ListItemOrTitleWidthPad)
		popUp.CastLocationOptionList.SetHeight(maxHeight - ListItemOrTitleWidthPad)
		return style.BorderStyle.
			MaxWidth(maxWidth).
			MaxHeight(maxHeight).
			Render(popUp.CastLocationOptionList.View())
	}
	return ""
}

// ----------------------------------
//
//	renders the cast-at-waypoint location option popup: sizes the list to
//	80 % width / 70 % height of the terminal, then wraps it in a rounded border;
//	returns an empty string if the popup model is the wrong type
//
// ----------------------------------
func renderCastWaypointLocationOptionPopUp(m *SpellbookInteractiveModel) string {
	popUp, ok := m.SpellPopUpModel.(*CastWaypointLocationOptionPopUpModel)
	if ok {
		maxWidth := int(float64(m.Width) * 0.8)
		maxHeight := int(float64(m.Height) * 0.7)
		popUp.CastWaypointLocationOptionList.SetWidth(maxWidth - ListItemOrTitleWidthPad)
		popUp.CastWaypointLocationOptionList.SetHeight(maxHeight - ListItemOrTitleWidthPad)
		return style.BorderStyle.
			MaxWidth(maxWidth).
			MaxHeight(maxHeight).
			Render(popUp.CastWaypointLocationOptionList.View())
	}
	return ""
}
