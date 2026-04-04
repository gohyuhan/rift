package ui

import (
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
