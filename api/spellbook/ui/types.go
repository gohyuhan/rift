package ui

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/x/ansi"
	"github.com/gohyuhan/rift/style"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
)

// HelpListItem holds the display data for a single row in the help popup.
type HelpListItem struct {
	Keybinding  string
	Action      string
	Description string
}

// ---------------------------------
//
//	list item and delegate types for the spell interactive UI
//
// ---------------------------------
type (
	spellInfoDelegate struct{}
	spellInfoItem     struct {
		SpellName    string
		SpellCommand []string
	}
)

// ----------------------------------
//
//	returns the value used when the list filters items; spells are
//	matched by name
//
// ----------------------------------
func (i spellInfoItem) FilterValue() string {
	return i.SpellName
}

// ----------------------------------
//
//	Height and Spacing define the row layout for the bubbles list delegate;
//	each spell occupies 2 lines (name + command ) with no
//	extra spacing; Update is a no-op as item-level updates are handled by
//	the parent model
//
// ----------------------------------
func (d spellInfoDelegate) Height() int                             { return 3 }
func (d spellInfoDelegate) Spacing() int                            { return 0 }
func (d spellInfoDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

// ----------------------------------
//
//	renders a single spell row as two lines: the spell name on the
//	first line and the command (indented) on the second; sealed spells are
//	rendered in a muted colour with a sealed label appended to the name,
//	while active spells use the vibrant purple palette; the selected
//	row is prefixed with a ❯ cursor, all others with two spaces
//
// ----------------------------------
func (d spellInfoDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(spellInfoItem)
	if !ok {
		return
	}

	componentWidth := m.Width() - ListItemOrTitleWidthPad

	spellName := fmt.Sprintf(" %s", i.SpellName)
	spellName = ansi.Truncate(spellName, componentWidth, "…")

	// command is indented with an extra space relative to the name
	spellCommand := fmt.Sprintf("   %s", strings.Join(i.SpellCommand, " "))
	spellCommand = ansi.Truncate(spellCommand, componentWidth, "…")

	spellName = style.RenderStringWithColor(spellName, style.ColorPurpleVibrant, false)
	spellCommand = style.RenderStringWithColor(spellCommand, style.ColorPurpleSoft, true)

	str := fmt.Sprintf("%s\n%s\n", spellName, spellCommand)

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
