package style

import (
	"image/color"

	"charm.land/lipgloss/v2"
)

var (
	ColorBlueSoft      = lipgloss.Color("#82AAFF")
	ColorBlueMuted     = lipgloss.Color("#A5B7E8")
	ColorYellowWarm    = lipgloss.Color("#F0D278")
	ColorYellowSoft    = lipgloss.Color("#F5E6A3")
	ColorGreenSoft     = lipgloss.Color("#98FB98")
	ColorError         = lipgloss.Color("#FF6B6B")
	ColorBlueVeryLight = lipgloss.Color("#E8F0FF")
	ColorBlueGrayMuted = lipgloss.Color("#6B7A9E")
	ColorPurpleSoft    = lipgloss.Color("#B496FF")
	ColorPurpleVibrant = lipgloss.Color("#9F7AEA")
	ColorCyanSoft      = lipgloss.Color("#7DD3FC")
)

// ----------------------------------
//
//	Renders text with the given foreground color. If faint is true, the text
//	is rendered with reduced intensity.
//
// ----------------------------------
func RenderStringWithColor(text string, color color.Color, faint bool) string {
	style := lipgloss.NewStyle().Foreground(color).Faint(faint)
	return style.Render(text)
}
