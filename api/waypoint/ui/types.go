package ui

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/x/ansi"
	"github.com/gohyuhan/rift/i18n"
	"github.com/gohyuhan/rift/style"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
)

// ---------------------------------
//
// for list component of git branch
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

func (i waypointInfoItem) FilterValue() string {
	return i.WaypointName
}

// for list component of Git branch
func (d waypointInfoDelegate) Height() int                             { return 1 }
func (d waypointInfoDelegate) Spacing() int                            { return 0 }
func (d waypointInfoDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d waypointInfoDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(waypointInfoItem)
	if !ok {
		return
	}

	componentWidth := m.Width() - ListItemOrTitleWidthPad

	waypointName := fmt.Sprintf(" %s", i.WaypointName)
	if i.WaypointIsSealed {
		waypointName = fmt.Sprintf(" %s %s", i.WaypointName, i18n.LANGUAGEMAPPING.RiftWaypointSealedLabel)
	}
	waypointName = ansi.Truncate(waypointName, componentWidth, "…")

	waypointPath := fmt.Sprintf("   %s", i.WaypointPath)
	waypointPath = ansi.Truncate(waypointPath, componentWidth, "…")

	if i.WaypointIsSealed {
		waypointName = style.RenderStringWithColor(waypointName, style.ColorSealedMuted, true)
		waypointPath = style.RenderStringWithColor(waypointPath, style.ColorSealedMuted, true)
	} else {
		waypointName = style.RenderStringWithColor(waypointName, style.ColorPurpleVibrant, false)
		waypointPath = style.RenderStringWithColor(waypointPath, style.ColorPurpleSoft, true)
	}

	str := fmt.Sprintf("%s\n%s", waypointName, waypointPath)

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
