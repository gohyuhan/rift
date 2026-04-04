package ui

import "charm.land/lipgloss/v2"

// ----------------------------------
//
//	composites the spell list and any active popup into a single string;
//	when a popup is showing, it is centered over the list using lipgloss layers;
//	when no popup is active the raw list view is returned directly
//
// ----------------------------------
func renderSpellbookInteractiveUIView(m *SpellbookInteractiveModel) string {
	if m.ShowPopUp.Load() {
		// Render the popup view into a string.
		popUpComponent := renderPopUpComponent(m)

		// Calculate the X and Y coordinates to center the popup.
		popUpWidth := lipgloss.Width(popUpComponent)
		popUpHeight := lipgloss.Height(popUpComponent)
		x := (m.Width - popUpWidth) / 2
		y := (m.Height - popUpHeight) / 2

		layers := []*lipgloss.Layer{
			lipgloss.NewLayer(m.SpellInfoList.View()),
			lipgloss.NewLayer(popUpComponent).X(x).Y(y).Z(1),
		}

		compositor := lipgloss.NewCompositor(layers...)

		return compositor.Render()
	}
	return m.SpellInfoList.View()
}
