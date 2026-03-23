package ui

import (
	"fmt"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
	"github.com/charmbracelet/x/ansi"
	apiUtils "github.com/gohyuhan/rift/api/utils"
	"github.com/gohyuhan/rift/db"
	"github.com/gohyuhan/rift/i18n"
	pb "github.com/gohyuhan/rift/proto"
	"github.com/gohyuhan/rift/style"
	"github.com/gohyuhan/rift/utils"
	"go.etcd.io/bbolt"
	"google.golang.org/protobuf/proto"
)

// ----------------------------------
//
//	reads every entry in the waypoint bucket and returns a slice of
//	waypointInfo records; uses a read-only View transaction so any
//	corruption writes are deferred to a separate Update after View
//	completes; corrupted proto entries are collected and recorded via
//	RecordCorruptedWaypointInfo before the error is returned to the caller
//
// ----------------------------------
func getAllWaypointsInfo(bboltDb *bbolt.DB) ([]waypointInfo, error) {
	var waypointsInfo []waypointInfo
	var corruptedWaypointName []string
	waypointCorrupted := false

	// seal updates are collected during the read-only View and applied afterwards;
	// calling a write transaction (Update) inside a View callback deadlocks bbolt
	type pendingSeal struct {
		name   string
		reason string
	}
	var toSeal []pendingSeal

	viewErr := bboltDb.View(func(tx *bbolt.Tx) error {
		// ensure the waypoint bucket exists before iterating
		waypointBucket := tx.Bucket(db.WaypointBucket)
		if waypointBucket == nil {
			return fmt.Errorf("%s", style.RenderStringWithColor(i18n.LANGUAGEMAPPING.WaypointBucketNotFoundError, style.ColorError, false))
		}

		// walk every key-value pair in the bucket
		retrieveError := waypointBucket.ForEach(func(k, v []byte) error {
			// deserialize the stored proto; capture the name, set the flag, and
			// return a sentinel error to stop ForEach — recording is deferred to
			// a separate Update after the View transaction completes
			existingWaypoint := &pb.Waypoint{}
			protoErr := proto.Unmarshal(v, existingWaypoint)

			// skip corrupted data
			if protoErr != nil {
				waypointCorrupted = true
				corruptedWaypointName = append(corruptedWaypointName, string(k))
				return nil
			}

			// verify the path still exists on disk; if not, mark for sealing after View closes
			isPathExist, isPathExistErr := utils.CheckIsPathExist(existingWaypoint.WaypointPath)
			if !isPathExist {
				existingWaypoint.WaypointIsSealed = true
				existingWaypoint.WaypointSealedReason = isPathExistErr.Error()
				toSeal = append(toSeal, pendingSeal{name: string(k), reason: existingWaypoint.WaypointSealedReason})
			}

			// construct the waypoint info type
			info := waypointInfo{
				WaypointName:         string(k),
				WaypointPath:         existingWaypoint.WaypointPath,
				WaypointIsSealed:     existingWaypoint.WaypointIsSealed,
				WaypointSealedReason: existingWaypoint.WaypointSealedReason,
			}

			waypointsInfo = append(waypointsInfo, info)

			return nil
		})

		// wrap ForEach failure into a user-facing message
		if retrieveError != nil {
			return fmt.Errorf("%s", style.RenderStringWithColor(i18n.LANGUAGEMAPPING.RiftWaypointRetrieveAllError, style.ColorError, false))
		}

		return nil
	})

	for _, s := range toSeal {
		apiUtils.UpdateWaypointIsSeal(bboltDb, s.name, true, s.reason)
	}

	if waypointCorrupted {
		viewErr = apiUtils.RecordCorruptedWaypointInfo(bboltDb, corruptedWaypointName)
	}

	return waypointsInfo, viewErr
}

// ----------------------------------
//
//	builds or rebuilds the bubbles list model on the WaypointInteractiveModel;
//	preserves the previously selected waypoint by name when re-initialising,
//	and falls back to the stored cursor position when no match is found;
//	called once after the first WindowSizeMsg so that layout dimensions
//	are valid before the list is rendered
//
// ----------------------------------
func initWaypointInfoListModel(m *WaypointInteractiveModel) error {
	previousSelectedWaypoint := m.WaypointInfoList.SelectedItem()
	selectedWayPointCursorPosition := -1

	latestWaypointInfoArray := []list.Item{}

	titleWidthLimit := m.Width - ListItemOrTitleWidthPad - ListTitleHorizontalPadding

	allWaypointsInfo, err := getAllWaypointsInfo(m.BboltDb)
	if err != nil {
		return err
	}

	if previousSelectedWaypoint != nil {
		// a previous selection exists; scan the fresh list to find a name match
		// so the cursor lands on the same waypoint after re-initialisation
		previousSelectedWaypointInfo := previousSelectedWaypoint.(waypointInfoItem)
		for index, waypoint := range allWaypointsInfo {
			if waypoint.WaypointName == previousSelectedWaypointInfo.WaypointName {
				selectedWayPointCursorPosition = index
			}
			latestWaypointInfoArray = append(latestWaypointInfoArray, waypointInfoItem(waypoint))
		}
	} else {
		// no previous selection; convert all records to list items
		for _, waypoint := range allWaypointsInfo {
			latestWaypointInfoArray = append(latestWaypointInfoArray, waypointInfoItem(waypoint))
		}
	}

	m.WaypointInfoList = list.New(latestWaypointInfoArray, waypointInfoDelegate{}, m.Width, m.Height)
	m.WaypointInfoList.SetShowPagination(false)
	m.WaypointInfoList.SetShowStatusBar(false)
	m.WaypointInfoList.SetFilteringEnabled(false)
	m.WaypointInfoList.SetShowFilter(false)
	m.WaypointInfoList.SetShowHelp(true)

	// truncate the title to prevent overflow when the terminal is narrow
	m.WaypointInfoList.Title = ansi.Truncate(i18n.LANGUAGEMAPPING.WaypointInfoListTitle, titleWidthLimit, "...")
	m.WaypointInfoList.Styles.Title = style.NewStyle.Bold(true)
	m.WaypointInfoList.Styles.PaginationStyle = style.NewStyle
	m.WaypointInfoList.Styles.TitleBar = style.NewStyle
	m.WaypointInfoList.Styles.HelpStyle = style.NewStyle
	m.WaypointInfoList.Help.Styles.ShortKey = style.NewStyle.Foreground(style.ColorBlueMuted).Bold(true)
	m.WaypointInfoList.Help.Styles.ShortDesc = style.NewStyle.Foreground(style.ColorBlueMuted)
	m.WaypointInfoList.Help.Styles.ShortSeparator = style.NewStyle.Foreground(style.ColorBlueGrayMuted)
	// clear the default key map so only our custom bindings appear in the help bar
	m.WaypointInfoList.KeyMap = list.KeyMap{}
	m.WaypointInfoList.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(key.WithKeys("q", "esc", "ctrl+c"), key.WithHelp("q/esc/ctrl+c", i18n.LANGUAGEMAPPING.ListQuitKeyHelp)),
		}
	}

	if selectedWayPointCursorPosition >= 0 {
		// restore the cursor to the previously selected waypoint by name match
		m.WaypointInfoList.Select(selectedWayPointCursorPosition)
		m.WaypointInfoListCursorPosition = selectedWayPointCursorPosition
	} else {
		// clamp the stored position to the last item if the list has shrunk
		if m.WaypointInfoListCursorPosition > len(m.WaypointInfoList.Items())-1 {
			m.WaypointInfoList.Select(len(m.WaypointInfoList.Items()) - 1)
			m.WaypointInfoListCursorPosition = len(m.WaypointInfoList.Items()) - 1
		} else {
			m.WaypointInfoList.Select(m.WaypointInfoListCursorPosition)
		}
	}

	// update the help key map to reflect whether the initial selection is sealed
	currentSelectedWaypoint := m.WaypointInfoList.SelectedItem()
	if currentSelectedWaypoint != nil {
		currentSelectedWaypointInfo := currentSelectedWaypoint.(waypointInfoItem)
		m.WaypointInfoList.AdditionalShortHelpKeys = initWaypointInfoListKeyMap(currentSelectedWaypointInfo.WaypointIsSealed)
	}

	return nil
}

// ----------------------------------
//
//	returns an AdditionalShortHelpKeys func tailored to the sealed state
//	of the currently selected waypoint; sealed items replace the enter
//	(navigate) binding with a u/U (unseal) binding, since navigation is
//	not permitted until the waypoint is explicitly unsealed
//
// ----------------------------------
func initWaypointInfoListKeyMap(isSealed bool) func() []key.Binding {
	if isSealed {
		return func() []key.Binding {
			return []key.Binding{
				key.NewBinding(key.WithKeys("↑", "k"), key.WithHelp("↑/k", i18n.LANGUAGEMAPPING.ListUpKeyHelp)),
				key.NewBinding(key.WithKeys("↓", "j"), key.WithHelp("↓/j", i18n.LANGUAGEMAPPING.ListDownKeyHelp)),
				key.NewBinding(key.WithKeys("u", "U"), key.WithHelp("u/U", i18n.LANGUAGEMAPPING.WaypointUnsealKeyHelp)),
				key.NewBinding(key.WithKeys("y"), key.WithHelp("y", i18n.LANGUAGEMAPPING.WaypointNameCopyKeyHelp)),
				key.NewBinding(key.WithKeys("Y"), key.WithHelp("Y", i18n.LANGUAGEMAPPING.WaypointPathCopyKeyHelp)),
				key.NewBinding(key.WithKeys("backspace"), key.WithHelp("backspace", i18n.LANGUAGEMAPPING.WaypointDestroyKeyHelp)),
				key.NewBinding(key.WithKeys("q", "esc", "ctrl+c"), key.WithHelp("q/esc/ctrl+c", i18n.LANGUAGEMAPPING.ListQuitKeyHelp)),
			}
		}
	} else {
		return func() []key.Binding {
			return []key.Binding{
				key.NewBinding(key.WithKeys("↑", "k"), key.WithHelp("↑/k", i18n.LANGUAGEMAPPING.ListUpKeyHelp)),
				key.NewBinding(key.WithKeys("↓", "j"), key.WithHelp("↓/j", i18n.LANGUAGEMAPPING.ListDownKeyHelp)),
				key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", i18n.LANGUAGEMAPPING.WaypointNavigateKeyHelp)),
				key.NewBinding(key.WithKeys("y"), key.WithHelp("y", i18n.LANGUAGEMAPPING.WaypointNameCopyKeyHelp)),
				key.NewBinding(key.WithKeys("Y"), key.WithHelp("Y", i18n.LANGUAGEMAPPING.WaypointPathCopyKeyHelp)),
				key.NewBinding(key.WithKeys("backspace"), key.WithHelp("backspace", i18n.LANGUAGEMAPPING.WaypointDestroyKeyHelp)),
				key.NewBinding(key.WithKeys("q", "esc", "ctrl+c"), key.WithHelp("q/esc/ctrl+c", i18n.LANGUAGEMAPPING.ListQuitKeyHelp)),
			}
		}
	}
}
