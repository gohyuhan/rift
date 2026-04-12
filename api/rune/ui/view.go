package ui

import "charm.land/lipgloss/v2"

// ----------------------------------
//
//	composites the rune engraving UI and any active popup into a single string;
//	when a popup is showing, it is centered over a blank background using lipgloss layers;
//	when no popup is active an empty string is returned
//
// ----------------------------------
func renderRuneInteractiveUIView(m *RuneInteractiveModel) string {
	if m.ShowPopUp.Load() {
		// Render the popup view into a string.
		popUpComponent := renderPopUpComponent(m)

		// Calculate the X and Y coordinates to center the popup.
		popUpWidth := lipgloss.Width(popUpComponent)
		popUpHeight := lipgloss.Height(popUpComponent)
		x := (m.Width - popUpWidth) / 2
		y := (m.Height - popUpHeight) / 2

		background := lipgloss.NewStyle().
			Width(m.Width).
			Height(m.Height).
			Render("")

		layers := []*lipgloss.Layer{
			lipgloss.NewLayer(background),
			lipgloss.NewLayer(popUpComponent).X(x).Y(y).Z(1),
		}

		compositor := lipgloss.NewCompositor(layers...)

		return compositor.Render()
	}
	return ""
}
