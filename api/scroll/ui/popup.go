package ui

import (
	"image/color"

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
func renderPopUpComponent(m *ScrollInteractiveModel) string {
	switch m.PopUpType {
	case HelpPopUp:
		return renderHelpPopUp(m)
	case InscribePopUp:
		return renderInscribePopUp(m)
	case InvokeLocationOptionPopUp:
		return renderInvokeLocationOptionPopUp(m)
	case InvokeWaypointLocationOptionPopUp:
		return renderInvokeWaypointLocationOptionPopUp(m)
	}
	return ""
}

// ----------------------------------
//
//	renders the help popup: sizes the viewport to 80 % width / 70 % height of
//	the terminal, then wraps it in a rounded border
//
// ----------------------------------
func renderHelpPopUp(m *ScrollInteractiveModel) string {
	maxWidth := int(float64(m.Width) * 0.8)
	maxHeight := int(float64(m.Height) * 0.7)
	m.RitualHelpViewport.SetWidth(maxWidth - ListItemOrTitleWidthPad)
	m.RitualHelpViewport.SetHeight(maxHeight - ListItemOrTitleWidthPad)
	return style.BorderStyle.
		MaxWidth(maxWidth).
		MaxHeight(maxHeight).
		Render(m.RitualHelpViewport.View())
}

// ----------------------------------
//
//	renders the inscribe new ritual popup: invokes RitualPopUpModel to *InscribePopUpModel
//	and sizes both text inputs to 80 % of the terminal width; joins the ritual
//	name title, ritual name input, command title, and command input vertically —
//	appending an inline error message below the inputs when one is present;
//	wraps the result in a rounded border sized to 80 % width / 70 % height;
//	returns an empty string if the model type assertion fails
//
// ----------------------------------
func renderInscribePopUp(m *ScrollInteractiveModel) string {
	popUp, ok := m.RitualPopUpModel.(*InscribePopUpModel)
	if ok {
		maxWidth := int(float64(m.Width) * 0.8)
		maxHeight := int(float64(m.Height) * 0.7)
		popUp.RitualNameInput.SetWidth(maxWidth - TextInputWidthPad)
		popUp.RitualDescriptionInput.SetWidth(maxWidth - TextInputWidthPad)
		popUp.RitualCommandsInput.SetWidth(maxWidth - TextInputWidthPad)
		strStyle := lipgloss.NewStyle().Width(maxWidth - TextInputWidthPad)
		var content string

		disabled := popUp.InscribeDisable.Load()
		buttonSelected := !popUp.RitualCommandsInput.Focused() && !popUp.RitualDescriptionInput.Focused() && !popUp.RitualNameInput.Focused()
		var buttonBorderColor color.Color
		if buttonSelected {
			buttonBorderColor = style.ColorBlueSoft
		} else {
			buttonBorderColor = style.ColorBlueGrayMuted
		}

		buttonLabel := style.RenderStringWithColor(i18n.LANGUAGEMAPPING.EngraveRuneEngraveButton, style.ColorBlueSoft, false)
		button := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(buttonBorderColor).
			Faint(disabled).
			Padding(0, 1).
			Render(buttonLabel)

		if popUp.EditMode {
			if popUp.Error != nil {
				errMessage := strStyle.Render(style.RenderStringWithColor(popUp.Error.Error(), style.ColorError, true))
				content = lipgloss.JoinVertical(
					lipgloss.Left,
					strStyle.Render(i18n.LANGUAGEMAPPING.RitualDescriptionInputTitle),
					popUp.RitualDescriptionInput.View(),
					strStyle.Render(i18n.LANGUAGEMAPPING.RitualCommandsInputTitle),
					popUp.RitualCommandsInput.View(),
					errMessage,
					button,
				)
			} else {
				content = lipgloss.JoinVertical(
					lipgloss.Left,
					strStyle.Render(i18n.LANGUAGEMAPPING.RitualDescriptionInputTitle),
					popUp.RitualDescriptionInput.View(),
					strStyle.Render(i18n.LANGUAGEMAPPING.RitualCommandsInputTitle),
					popUp.RitualCommandsInput.View(),
					button,
				)
			}
		} else {
			if popUp.Error != nil {
				errMessage := strStyle.Render(style.RenderStringWithColor(popUp.Error.Error(), style.ColorError, true))
				content = lipgloss.JoinVertical(
					lipgloss.Left,
					strStyle.Render(i18n.LANGUAGEMAPPING.RitualNameInputTitle),
					popUp.RitualNameInput.View(),
					strStyle.Render(i18n.LANGUAGEMAPPING.RitualDescriptionInputTitle),
					popUp.RitualDescriptionInput.View(),
					strStyle.Render(i18n.LANGUAGEMAPPING.RitualCommandsInputTitle),
					popUp.RitualCommandsInput.View(),
					errMessage,
					button,
				)
			} else {
				content = lipgloss.JoinVertical(
					lipgloss.Left,
					strStyle.Render(i18n.LANGUAGEMAPPING.RitualNameInputTitle),
					popUp.RitualNameInput.View(),
					strStyle.Render(i18n.LANGUAGEMAPPING.RitualDescriptionInputTitle),
					popUp.RitualDescriptionInput.View(),
					strStyle.Render(i18n.LANGUAGEMAPPING.RitualCommandsInputTitle),
					popUp.RitualCommandsInput.View(),
					button,
				)
			}
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
//	renders the invoke-location option popup: sizes the list to 80 % of the
//	terminal width and a fixed 12-row height, then wraps it in a rounded border;
//	returns an empty string if the popup model is the wrong type
//
// ----------------------------------
func renderInvokeLocationOptionPopUp(m *ScrollInteractiveModel) string {
	popUp, ok := m.RitualPopUpModel.(*InvokeLocationOptionPopUpModel)
	if ok {
		maxWidth := int(float64(m.Width) * 0.8)
		maxHeight := 12
		popUp.InvokeLocationOptionList.SetWidth(maxWidth - ListItemOrTitleWidthPad)
		popUp.InvokeLocationOptionList.SetHeight(maxHeight - ListItemOrTitleWidthPad)
		return style.BorderStyle.
			MaxWidth(maxWidth).
			MaxHeight(maxHeight).
			Render(popUp.InvokeLocationOptionList.View())
	}
	return ""
}

// ----------------------------------
//
//	renders the invoke-at-waypoint location option popup: sizes the list to
//	80 % width / 70 % height of the terminal, then wraps it in a rounded border;
//	returns an empty string if the popup model is the wrong type
//
// ----------------------------------
func renderInvokeWaypointLocationOptionPopUp(m *ScrollInteractiveModel) string {
	popUp, ok := m.RitualPopUpModel.(*InvokeWaypointLocationOptionPopUpModel)
	if ok {
		maxWidth := int(float64(m.Width) * 0.8)
		maxHeight := int(float64(m.Height) * 0.7)
		popUp.InvokeWaypointLocationOptionList.SetWidth(maxWidth - ListItemOrTitleWidthPad)
		popUp.InvokeWaypointLocationOptionList.SetHeight(maxHeight - ListItemOrTitleWidthPad)
		return style.BorderStyle.
			MaxWidth(maxWidth).
			MaxHeight(maxHeight).
			Render(popUp.InvokeWaypointLocationOptionList.View())
	}
	return ""
}
