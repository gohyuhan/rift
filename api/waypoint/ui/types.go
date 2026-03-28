package ui

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/x/ansi"
	"github.com/gohyuhan/rift/i18n"
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
//	list item and delegate types for the waypoint interactive UI
//
// ---------------------------------
type (
	waypointInfoDelegate struct{}
	waypointInfoItem     struct {
		WaypointName         string
		WaypointPath         string
		WaypointIsSealed     bool
		WaypointSealedReason string
	}
)

// ----------------------------------
//
//	returns the value used when the list filters items; waypoints are
//	matched by name
//
// ----------------------------------
func (i waypointInfoItem) FilterValue() string {
	return i.WaypointName
}

// ----------------------------------
//
//	Height and Spacing define the row layout for the bubbles list delegate;
//	each waypoint occupies 3 lines (name + path + sealed reason) with no
//	extra spacing; the sealed-reason line is empty for active waypoints;
//	Update is a no-op as item-level updates are handled by the parent model
//
// ----------------------------------
func (d waypointInfoDelegate) Height() int                             { return 3 }
func (d waypointInfoDelegate) Spacing() int                            { return 0 }
func (d waypointInfoDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

// ----------------------------------
//
//	renders a single waypoint row as two lines: the waypoint name on the
//	first line and the path (indented) on the second; sealed waypoints are
//	rendered in a muted colour with a sealed label appended to the name,
//	while active waypoints use the vibrant purple palette; the selected
//	row is prefixed with a ❯ cursor, all others with two spaces
//
// ----------------------------------
func (d waypointInfoDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(waypointInfoItem)
	if !ok {
		return
	}

	componentWidth := m.Width() - ListItemOrTitleWidthPad

	// append sealed label when the waypoint cannot be navigated to
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

// RebindPopUpModel holds the state for the rebind path input popup.
type RebindPopUpModel struct {
	RebindPathInput    textinput.Model
	WaypointName       string
	Error              error
	OnInputFuncTrigger func(bboltDb *bbolt.DB, waypointName string, rebindTo string, logToTerminal bool) error
}

// ReforgePopUpModel holds the state for the reforge name input popup.
type ReforgePopUpModel struct {
	ReforgeNameInput   textinput.Model
	WaypointName       string
	Error              error
	OnInputFuncTrigger func(bboltDb *bbolt.DB, waypointName string, reforgeTo string, logToTerminal bool) error
}
