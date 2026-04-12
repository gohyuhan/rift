package ui

import (
	"fmt"
	"io"
	"strings"
	"sync/atomic"

	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/textarea"
	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/gohyuhan/rift/style"
)

// ---------------------------------
//
//	list item and delegate types for the runeEngraveOption list
//
// ---------------------------------
type (
	runeEngraveOptionInfoDelegate struct{}
	runeEngraveOptionInfoItem     struct {
		runeEngraveOptionName string
		runeEngraveOptionDesc string
		runeEngraveOptionType string
	}
)

// ----------------------------------
//
//	returns the value used when the list filters items; runeEngraveOptions are
//	matched by name
//
// ----------------------------------
func (i runeEngraveOptionInfoItem) FilterValue() string {
	return i.runeEngraveOptionName
}

// ----------------------------------
//
//	Height and Spacing define the row layout for the bubbles list delegate;
//	each runeEngraveOption occupies 3 lines (name + command + blank separator) with no
//	extra spacing; Update is a no-op as item-level updates are handled by
//	the parent model
//
// ----------------------------------
func (d runeEngraveOptionInfoDelegate) Height() int                             { return 3 }
func (d runeEngraveOptionInfoDelegate) Spacing() int                            { return 0 }
func (d runeEngraveOptionInfoDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

// ----------------------------------
//
//	renders a single runeEngraveOption row as three lines: the runeEngraveOption name on the
//	first line, the command (indented) on the second, and a blank line
//	as a separator; the name is coloured with the vibrant purple palette
//	and the command with the soft purple palette; the selected row is
//	prefixed with a ❯ cursor, all others with two spaces
//
// ----------------------------------
func (d runeEngraveOptionInfoDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(runeEngraveOptionInfoItem)
	if !ok {
		return
	}

	componentWidth := m.Width() - ListItemOrTitleWidthPad

	runeEngraveOptionName := fmt.Sprintf(" %s", i.runeEngraveOptionName)
	runeEngraveOptionName = ansi.Truncate(runeEngraveOptionName, componentWidth, "…")

	// command is indented with an extra space relative to the name
	runeEngraveOptionCommand := fmt.Sprintf("   %s", i.runeEngraveOptionDesc)
	runeEngraveOptionCommand = ansi.Truncate(runeEngraveOptionCommand, componentWidth, "…")

	runeEngraveOptionName = style.RenderStringWithColor(runeEngraveOptionName, style.ColorPurpleVibrant, false)
	runeEngraveOptionCommand = style.RenderStringWithColor(runeEngraveOptionCommand, style.ColorPurpleSoft, true)

	str := fmt.Sprintf("%s\n%s\n", runeEngraveOptionName, runeEngraveOptionCommand)

	// prefix the selected row with a cursor glyph; all others get padding
	var fn func(...string) string
	if index == m.Index() {
		fn = func(s ...string) string {
			return style.SelectedItemStyle.Render("❯ " + strings.Join(s, " "))
		}
	} else {
		fn = func(s ...string) string {
			return style.ItemStyle.Render("  " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

type ChooseRuneEngraveOptionPopUpModel struct {
	EngraveOptionList list.Model
}

type EngraveRuneCommandsPopUpModel struct {
	RuneCommandsTextArea  textarea.Model
	RuneEngraveOptionType string
	TextAreaFocused       atomic.Bool
	EngraveDisable        atomic.Bool
	Error                 error
}
