package ui

import (
	"image/color"

	"charm.land/lipgloss/v2"
	apiUtils "github.com/gohyuhan/rift/api/utils"
	"github.com/gohyuhan/rift/i18n"
	"github.com/gohyuhan/rift/style"
)

// ----------------------------------
//
//	dispatches to the appropriate popup renderer based on the active PopUpType;
//	returns an empty string if the type is unrecognised or not yet implemented
//
// ----------------------------------
func renderPopUpComponent(m *RuneInteractiveModel) string {
	switch m.PopUpType {
	case ChooseRuneEngraveOptionPopUp:
		return renderChooseRuneEngraveOptionPopUp(m)
	case EngraveRuneCommandsPopUp:
		return renderEngraveRuneCommandsPopUp(m)
	}
	return ""
}

// ----------------------------------
//
//	renders the rune engrave/remove option list popup: sizes the list to 80 %
//	of the terminal width; wraps plain title strings with a width-constrained
//	style so they word-wrap instead of overflowing; then wraps the block in a
//	rounded border sized to 80 % width / 70 % height;
//	returns an empty string if the popup model is the wrong type
//
// ----------------------------------
func renderChooseRuneEngraveOptionPopUp(m *RuneInteractiveModel) string {
	popUp, ok := m.RunePopUpModel.(*ChooseRuneEngraveOptionPopUpModel)
	if ok {
		maxWidth := int(float64(m.Width) * 0.8)
		maxHeight := int(float64(m.Height) * 0.7)
		popUp.EngraveOptionList.SetWidth(maxWidth - ListItemOrTitleWidthPad)
		popUp.EngraveOptionList.SetHeight(len(popUp.EngraveOptionList.Items()) * 4)
		strStyle := lipgloss.NewStyle().Width(maxWidth - ListItemOrTitleWidthPad)
		content := lipgloss.JoinVertical(
			lipgloss.Left,
			strStyle.Render(i18n.LANGUAGEMAPPING.RuneEngraveTypeOptionListTitle),
			popUp.EngraveOptionList.View(),
		)
		return style.BorderStyle.
			MaxWidth(maxWidth).
			MaxHeight(maxHeight).
			Render(content)
	}
	return ""
}

// ----------------------------------
//
//	renders the rune command entry popup for the selected slot (enter or leave):
//	sizes the textarea to 65 % of the popup height; validates the textarea
//	content on every render and toggles the Engrave button's disabled state
//	accordingly; wraps plain title and error strings with a width-constrained
//	style so they word-wrap instead of overflowing; the Engrave button border
//	is highlighted when it has focus (tab/shift+tab toggles focus);
//	appends a coloured error line below the textarea when validation fails;
//	wraps the whole block in a rounded border sized to 80 % width / 70 % height;
//	returns an empty string if the popup model is the wrong type
//
// ----------------------------------
func renderEngraveRuneCommandsPopUp(m *RuneInteractiveModel) string {
	popUp, ok := m.RunePopUpModel.(*EngraveRuneCommandsPopUpModel)
	if ok {
		maxWidth := int(float64(m.Width) * 0.8)
		maxHeight := int(float64(m.Height) * 0.7)

		textAreaHeight := int(float64(maxHeight) * 0.65)
		popUp.RuneCommandsTextArea.SetWidth(maxWidth - TextInputWidthPad)
		popUp.RuneCommandsTextArea.SetHeight(textAreaHeight)
		var title string
		switch popUp.RuneEngraveOptionType {
		case EngraveRuneEnterType:
			title = i18n.LANGUAGEMAPPING.EngraveRuneEnterTitle
		case EngraveRuneLeaveType:
			title = i18n.LANGUAGEMAPPING.EngraveRuneLeaveTitle
		}

		normRuneArray, runeCmdsErr := apiUtils.NormalizeAndCheckRuneCommandsAreValid(popUp.RuneCommandsTextArea.Value())
		if runeCmdsErr != nil {
			popUp.Error = runeCmdsErr
			popUp.EngraveDisable.Store(true)
		} else {
			popUp.Error = nil
			popUp.EngraveDisable.Store(false)
		}

		if len(normRuneArray) > 0 {
			popUp.EngraveDisable.Store(false)
		} else {
			popUp.EngraveDisable.Store(true)
		}

		disabled := popUp.EngraveDisable.Load()
		textAreaFocused := popUp.TextAreaFocused.Load()
		buttonSelected := !textAreaFocused
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

		strStyle := lipgloss.NewStyle().Width(maxWidth - TextInputWidthPad)
		var content string
		if popUp.Error != nil {
			errMessage := strStyle.Render(style.RenderStringWithColor(popUp.Error.Error(), style.ColorError, true))
			content = lipgloss.JoinVertical(
				lipgloss.Left,
				strStyle.Render(title),
				popUp.RuneCommandsTextArea.View(),
				errMessage,
				button,
			)
		} else {
			content = lipgloss.JoinVertical(
				lipgloss.Left,
				strStyle.Render(title),
				popUp.RuneCommandsTextArea.View(),
				button,
			)
		}
		return style.BorderStyle.
			MaxWidth(maxWidth).
			MaxHeight(maxHeight).
			Render(content)
	}
	return ""
}
