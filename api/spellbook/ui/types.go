package ui

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/x/ansi"
	"github.com/gohyuhan/rift/style"
	"go.etcd.io/bbolt"

	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/textinput"
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
//	each spell occupies 3 lines (name + command + blank separator) with no
//	extra spacing; Update is a no-op as item-level updates are handled by
//	the parent model
//
// ----------------------------------
func (d spellInfoDelegate) Height() int                             { return 3 }
func (d spellInfoDelegate) Spacing() int                            { return 0 }
func (d spellInfoDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

// ----------------------------------
//
//	renders a single spell row as three lines: the spell name on the
//	first line, the command (indented) on the second, and a blank line
//	as a separator; the name is coloured with the vibrant purple palette
//	and the command with the soft purple palette; the selected row is
//	prefixed with a ❯ cursor, all others with two spaces
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

// ----------------------------------
//
//	LearnPopUpModel holds the state for the learn new spell popup
//
// ----------------------------------
type LearnPopUpModel struct {
	SpellNameInput         textinput.Model
	SpellCommandInput      textinput.Model
	TotalInputField        int
	CurrentFocusInputIndex int
	Error                  error
	OnInputFuncTrigger     func(bboltDb *bbolt.DB, spellName string, spellCommand string) (bool, error)
}

type CastLocationOptionPopUpModel struct {
	CastLocationOptionList list.Model
}

// ---------------------------------
//
//	list item and delegate types for the cast location option
//
// ---------------------------------
type (
	castLocationOptionDelegate struct{}
	castLocationOptionItem     struct {
		Title       string
		Description string
		OptionType  string
	}
)

// ----------------------------------
//
//	returns the value used when the list filters items; cast location
//	options are matched by title
//
// ----------------------------------
func (i castLocationOptionItem) FilterValue() string {
	return i.Title
}

// ----------------------------------
//
//	Height and Spacing define the row layout for the bubbles list delegate;
//	each option occupies 3 lines (title + description + blank separator)
//	with no extra spacing; Update is a no-op as item-level updates are
//	handled by the parent model
//
// ----------------------------------
func (d castLocationOptionDelegate) Height() int                             { return 3 }
func (d castLocationOptionDelegate) Spacing() int                            { return 0 }
func (d castLocationOptionDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

// ----------------------------------
//
//	renders a single cast location option row as three lines: the option
//	title on the first line, the description (indented) on the second, and
//	a blank line as a separator; the title is coloured with the vibrant
//	purple palette and the description with the soft purple palette; the
//	selected row is prefixed with a ❯ cursor, all others with two spaces
//
// ----------------------------------
func (d castLocationOptionDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(castLocationOptionItem)
	if !ok {
		return
	}

	componentWidth := m.Width() - ListItemOrTitleWidthPad

	optionTitle := fmt.Sprintf(" %s", i.Title)
	optionTitle = ansi.Truncate(optionTitle, componentWidth, "…")

	// description is indented with an extra space relative to the title
	optionDesc := fmt.Sprintf("   %s", i.Description)
	optionDesc = ansi.Truncate(optionDesc, componentWidth, "…")

	optionTitle = style.RenderStringWithColor(optionTitle, style.ColorPurpleVibrant, false)
	optionDesc = style.RenderStringWithColor(optionDesc, style.ColorPurpleSoft, true)

	str := fmt.Sprintf("%s\n%s\n", optionTitle, optionDesc)

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
