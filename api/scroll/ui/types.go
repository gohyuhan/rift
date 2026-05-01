package ui

import (
	"fmt"
	"io"
	"strings"
	"sync/atomic"

	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/gohyuhan/rift/i18n"
	pb "github.com/gohyuhan/rift/proto"
	"github.com/gohyuhan/rift/style"
)

// HelpListItem holds the display data for a single row in the help popup.
type HelpListItem struct {
	Keybinding  string
	Action      string
	Description string
}

// ---------------------------------
//
//	list item and delegate types for the ritual interactive UI
//
// ---------------------------------
type (
	ritualInfoDelegate struct{}
	ritualInfoItem     struct {
		RitualName string
		RitualDesc string
	}
)

// ----------------------------------
//
//	returns the value used when the list filters items; rituals are
//	matched by name
//
// ----------------------------------
func (i ritualInfoItem) FilterValue() string {
	return i.RitualName
}

// ----------------------------------
//
//	Height and Spacing define the row layout for the bubbles list delegate;
//	each ritual occupies 3 lines (name + command + blank separator) with no
//	extra spacing; Update is a no-op as item-level updates are handled by
//	the parent model
//
// ----------------------------------
func (d ritualInfoDelegate) Height() int                             { return 4 }
func (d ritualInfoDelegate) Spacing() int                            { return 0 }
func (d ritualInfoDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

// ----------------------------------
//
//	renders a single ritual row as three lines: the ritual name on the
//	first line, the command (indented) on the second, and a blank line
//	as a separator; the name is coloured with the vibrant purple palette
//	and the command with the soft purple palette; the selected row is
//	prefixed with a ❯ cursor, all others with two spaces
//
// ----------------------------------
func (d ritualInfoDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(ritualInfoItem)
	if !ok {
		return
	}

	componentWidth := m.Width() - ListItemOrTitleWidthPad

	ritualName := fmt.Sprintf(" %s", i.RitualName)
	ritualName = ansi.Truncate(ritualName, componentWidth, "…")

	// command is indented with an extra space relative to the name
	parsedRitualDescArray := strings.Split(i.RitualDesc, "\n")
	var parsedRitualDesc strings.Builder
	for index, ritualDesc := range parsedRitualDescArray {
		if index > 1 {
			parsedRitualDesc.WriteString("...")
			break
		} else {
			parsedRitualDesc.WriteString("\n")
		}
		parsedRitualDesc.WriteString("   " + ritualDesc)
	}
	ritualDesc := fmt.Sprintf("%s", parsedRitualDesc.String())
	ritualDesc = ansi.Truncate(ritualDesc, componentWidth, "…")

	ritualName = style.RenderStringWithColor(ritualName, style.ColorPurpleVibrant, false)
	ritualDesc = style.RenderStringWithColor(ritualDesc, style.ColorPurpleSoft, true)

	str := fmt.Sprintf("%s\n%s\n", ritualName, ritualDesc)

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
//	InscribePopUpModel holds the state for the learn new ritual popup
//
// ----------------------------------
type InscribePopUpModel struct {
	RitualName             string
	RitualNameInput        textinput.Model
	RitualDescriptionInput textarea.Model
	RitualCommandsInput    textarea.Model
	EditMode               bool // indicate whether this is for edit existing ritual or inscribe new ritual, true for edit, false for inscribe new ritual
	TotalInputField        int
	CurrentFocusInputIndex int
	Error                  error
	InscribeDisable        atomic.Bool
	OnInputFuncTrigger     func(ritualName string, ritualDesc string, ritualCmdString string, override bool) error
}

type InvokeLocationOptionPopUpModel struct {
	InvokeLocationOptionList list.Model
	SelectedRitualName       string
}

// ---------------------------------
//
//	list item and delegate types for the invoke location option
//
// ---------------------------------
type (
	invokeLocationOptionDelegate struct{}
	invokeLocationOptionItem     struct {
		Title       string
		Description string
		OptionType  string
	}
)

// ----------------------------------
//
//	returns the value used when the list filters items; invoke location
//	options are matched by title
//
// ----------------------------------
func (i invokeLocationOptionItem) FilterValue() string {
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
func (d invokeLocationOptionDelegate) Height() int                             { return 3 }
func (d invokeLocationOptionDelegate) Spacing() int                            { return 0 }
func (d invokeLocationOptionDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

// ----------------------------------
//
//	renders a single invoke location option row as three lines: the option
//	title on the first line, the description (indented) on the second, and
//	a blank line as a separator; the title is coloured with the vibrant
//	purple palette and the description with the soft purple palette; the
//	selected row is prefixed with a ❯ cursor, all others with two spaces
//
// ----------------------------------
func (d invokeLocationOptionDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(invokeLocationOptionItem)
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

type InvokeWaypointLocationOptionPopUpModel struct {
	InvokeWaypointLocationOptionList list.Model
	SelectedRitualName               string
}

// ---------------------------------
//
//	list item and delegate types for the invoke-to-waypoint location popup
//
// ---------------------------------
type (
	invokeWaypointLocationDelegate   struct{}
	invokeWaypointLocationOptionItem struct {
		WaypointName           string
		WaypointPath           string
		WaypointAddedAt        string
		WaypointTravelledCount int64
		WaypointIsSealed       bool
		WaypointSealedReason   string
		EnterRune              []*pb.RuneCmds
		LeaveRune              []*pb.RuneCmds
	}
)

// ----------------------------------
//
//	returns the value used when the list filters items; invoke-to-waypoint
//	options are matched by waypoint name
//
// ----------------------------------
func (i invokeWaypointLocationOptionItem) FilterValue() string {
	return i.WaypointName
}

// ----------------------------------
//
//	Height and Spacing define the row layout for the bubbles list delegate;
//	each waypoint option occupies 3 lines (name + path + sealed reason) with
//	no extra spacing; the sealed-reason line is blank for unsealed waypoints;
//	Update is a no-op as item-level updates are handled by the parent model
//
// ----------------------------------
func (d invokeWaypointLocationDelegate) Height() int                             { return 3 }
func (d invokeWaypointLocationDelegate) Spacing() int                            { return 0 }
func (d invokeWaypointLocationDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

// ----------------------------------
//
//	renders a single invoke-to-waypoint option as three lines: the waypoint
//	name on the first line, the path (indented) on the second, and the sealed
//	reason (indented) on the third; sealed waypoints are rendered in a muted
//	colour with a sealed label appended to the name, while active waypoints
//	use the vibrant purple palette; the selected row is prefixed with a ❯
//	cursor, all others with two spaces
//
// ----------------------------------
func (d invokeWaypointLocationDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(invokeWaypointLocationOptionItem)
	if !ok {
		return
	}

	componentWidth := m.Width() - ListItemOrTitleWidthPad

	// append sealed label when the waypoint cannot be invoke to
	waypointName := fmt.Sprintf(" %s", i.WaypointName)
	if i.WaypointIsSealed {
		waypointName = fmt.Sprintf(" %s %s", i.WaypointName, i18n.LANGUAGEMAPPING.RiftWaypointSealedLabel)
	}
	waypointName = ansi.Truncate(waypointName, componentWidth, "…")

	// path is indented with an extra space relative to the name
	waypointPath := fmt.Sprintf("   %s", i.WaypointPath)
	waypointPath = ansi.Truncate(waypointPath, componentWidth, "…")

	var waypointSealedReason string

	// apply colour based on sealed state
	if i.WaypointIsSealed {
		waypointName = style.RenderStringWithColor(waypointName, style.ColorSealedMuted, true)
		waypointPath = style.RenderStringWithColor(waypointPath, style.ColorSealedMuted, true)

		waypointSealedReason = fmt.Sprintf("   %s %s", i18n.LANGUAGEMAPPING.RiftWaypointDetailSealedReason, i.WaypointSealedReason)
		waypointSealedReason = ansi.Truncate(waypointSealedReason, componentWidth, "...")
		waypointSealedReason = style.RenderStringWithColor(waypointSealedReason, style.ColorSealedMuted, true)
	} else {
		waypointName = style.RenderStringWithColor(waypointName, style.ColorPurpleVibrant, false)
		waypointPath = style.RenderStringWithColor(waypointPath, style.ColorPurpleSoft, true)
	}

	str := fmt.Sprintf("%s\n%s\n%s", waypointName, waypointPath, waypointSealedReason)

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
