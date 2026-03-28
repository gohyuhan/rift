package style

import (
	"image/color"
	"strings"
	"unicode/utf8"

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
	ColorSealedMuted   = lipgloss.Color("#5C5470")
)

var (
	SelectedItemStyle = lipgloss.NewStyle().Foreground(ColorPurpleVibrant).Bold(true)
	ItemStyle         = lipgloss.NewStyle().Foreground(ColorPurpleSoft)
	NewStyle          = lipgloss.NewStyle()
	BorderStyle       = NewStyle.Border(lipgloss.RoundedBorder()).Padding(0).Margin(0).BorderForeground(ColorBlueGrayMuted)
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

// ----------------------------------
//
//	PadAndRenderLabels takes a slice of raw label strings, finds the longest
//	by rune count, pads all labels to that width, and returns each one
//	rendered with the given color and faint setting.
//
// ----------------------------------
func PadAndRenderLabels(labels []string, color color.Color, faint bool) []string {
	maxLen := 0
	for _, l := range labels {
		if n := utf8.RuneCountInString(l); n > maxLen {
			maxLen = n
		}
	}

	rendered := make([]string, len(labels))
	for i, l := range labels {
		padded := l + strings.Repeat(" ", maxLen-utf8.RuneCountInString(l))
		rendered[i] = RenderStringWithColor(padded, color, faint)
	}
	return rendered
}
