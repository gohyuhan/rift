package ui

import (
	"fmt"
	"strings"

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

	// best-effort: persist each seal to the DB; failures are silently ignored —
	// the in-memory waypointInfo records already carry WaypointIsSealed=true
	// so the UI reflects the correct sealed state regardless
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
			key.NewBinding(key.WithKeys("q", "esc", "ctrl+c"), key.WithHelp("q/esc/ctrl+c", i18n.LANGUAGEMAPPING.WaypointUIQuitKeyHelp)),
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
		m.WaypointInfoList.AdditionalShortHelpKeys = initShortWaypointInfoListKeyMap(NoPopUp, currentSelectedWaypointInfo.WaypointIsSealed)
	}

	return nil
}

// ----------------------------------
//
//	returns an AdditionalShortHelpKeys func whose bindings are determined by
//	popUpType and, when popUpType is NoPopUp, by isSealed; for NoPopUp it
//	shows ? (help) plus the primary action key relevant to the item's sealed
//	state, keeping the status bar compact; for RebindPopUp and ReforgePopUp
//	it shows only the popup-relevant keys: enter (submit), esc (close), and
//	q/ctrl+c (quit)
//
// ----------------------------------
func initShortWaypointInfoListKeyMap(popUpType string, isSealed bool) func() []key.Binding {
	switch popUpType {
	case NoPopUp:
		if isSealed {
			return func() []key.Binding {
				return []key.Binding{
					key.NewBinding(key.WithKeys("?"), key.WithHelp("?", i18n.LANGUAGEMAPPING.WaypointUIHelpKeyHelp)),
					key.NewBinding(key.WithKeys("q", "esc", "ctrl+c"), key.WithHelp("q/esc/ctrl+c", i18n.LANGUAGEMAPPING.WaypointUIQuitKeyHelp)),
				}
			}
		} else {
			return func() []key.Binding {
				return []key.Binding{
					key.NewBinding(key.WithKeys("?"), key.WithHelp("?", i18n.LANGUAGEMAPPING.WaypointUIHelpKeyHelp)),
					key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", i18n.LANGUAGEMAPPING.WaypointNavigateKeyHelp)),
					key.NewBinding(key.WithKeys("q", "esc", "ctrl+c"), key.WithHelp("q/esc/ctrl+c", i18n.LANGUAGEMAPPING.WaypointUIQuitKeyHelp)),
				}
			}
		}
	case RebindPopUp:
		return func() []key.Binding {
			return []key.Binding{
				key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", i18n.LANGUAGEMAPPING.WaypointRebindKeyHelp)),
				key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", i18n.LANGUAGEMAPPING.WaypointClosePopUp)),
				key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q/ctrl+c", i18n.LANGUAGEMAPPING.WaypointUIQuitKeyHelp)),
			}
		}
	case ReforgePopUp:
		return func() []key.Binding {
			return []key.Binding{
				key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", i18n.LANGUAGEMAPPING.WaypointReforgeKeyHelp)),
				key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", i18n.LANGUAGEMAPPING.WaypointClosePopUp)),
				key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q/ctrl+c", i18n.LANGUAGEMAPPING.WaypointUIQuitKeyHelp)),
			}
		}
	case HelpPopUp:
		return func() []key.Binding {
			return []key.Binding{
				key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", i18n.LANGUAGEMAPPING.WaypointClosePopUp)),
			}
		}
	}

	return func() []key.Binding { return []key.Binding{} }
}

// ----------------------------------
//
//	builds the full key binding reference and stores it as rendered text in
//	the help viewport; each entry is formatted as keybinding → action →
//	description, all coloured with the purple palette; called once during
//	model initialisation so the content is ready before the first render
//
// ----------------------------------
func initWaypointHelpViewport(m *WaypointInteractiveModel) {
	helpList := []HelpListItem{
		{Keybinding: "?", Action: i18n.LANGUAGEMAPPING.WaypointUIHelpKeyHelp, Description: i18n.LANGUAGEMAPPING.WaypointUIHelpKeyHelpDescription},
		{Keybinding: "↑/k", Action: i18n.LANGUAGEMAPPING.WaypointUIUpKeyHelp, Description: i18n.LANGUAGEMAPPING.WaypointUIUpKeyHelpDescription},
		{Keybinding: "↓/j", Action: i18n.LANGUAGEMAPPING.WaypointUIDownKeyHelp, Description: i18n.LANGUAGEMAPPING.WaypointUIDownKeyHelpDescription},
		{Keybinding: "enter", Action: i18n.LANGUAGEMAPPING.WaypointNavigateKeyHelp, Description: i18n.LANGUAGEMAPPING.WaypointNavigateKeyHelpDescription},
		{Keybinding: "u/U", Action: i18n.LANGUAGEMAPPING.WaypointUnsealKeyHelp, Description: i18n.LANGUAGEMAPPING.WaypointUnsealKeyHelpDescription},
		{Keybinding: "r", Action: i18n.LANGUAGEMAPPING.WaypointRebindKeyHelp, Description: i18n.LANGUAGEMAPPING.WaypointRebindKeyHelpDescription},
		{Keybinding: "R", Action: i18n.LANGUAGEMAPPING.WaypointReforgeKeyHelp, Description: i18n.LANGUAGEMAPPING.WaypointReforgeKeyHelpDescription},
		{Keybinding: "y/Y", Action: i18n.LANGUAGEMAPPING.WaypointNameCopyPathCopyKeyHelp, Description: i18n.LANGUAGEMAPPING.WaypointNameCopyPathCopyKeyHelpDescription},
		{Keybinding: "backspace", Action: i18n.LANGUAGEMAPPING.WaypointDestroyKeyHelp, Description: i18n.LANGUAGEMAPPING.WaypointDestroyKeyHelpDescription},
		{Keybinding: "ctrl+y", Action: i18n.LANGUAGEMAPPING.WaypointCopyFromInputValueKeyHelp, Description: i18n.LANGUAGEMAPPING.WaypointCopyFromInputValueKeyHelpDescription},
		{Keybinding: "ctrl+p", Action: i18n.LANGUAGEMAPPING.WaypointPasteIntoInputValueKeyHelp, Description: i18n.LANGUAGEMAPPING.WaypointPasteIntoInputValueKeyHelpDescription},
		{Keybinding: "esc", Action: i18n.LANGUAGEMAPPING.WaypointClosePopUp, Description: i18n.LANGUAGEMAPPING.WaypointClosePopUpDescription},
		{Keybinding: "q/esc/ctrl+c", Action: i18n.LANGUAGEMAPPING.WaypointUIQuitKeyHelp, Description: i18n.LANGUAGEMAPPING.WaypointUIQuitKeyHelpDescription},
	}

	var content strings.Builder
	for _, helpitem := range helpList {
		keybinding := style.RenderStringWithColor(fmt.Sprintf("[%s]", helpitem.Keybinding), style.ColorPurpleVibrant, false)
		action := style.RenderStringWithColor(helpitem.Action, style.ColorPurpleSoft, false)
		description := style.RenderStringWithColor(helpitem.Description, style.ColorPurpleSoft, true)
		content.WriteString(keybinding)
		content.WriteRune('\n')
		content.WriteString(action)
		content.WriteRune('\n')
		content.WriteString(description)
		content.WriteRune('\n')
		content.WriteRune('\n')
	}

	m.WaypointHelpViewport.SetContent(content.String())
}
