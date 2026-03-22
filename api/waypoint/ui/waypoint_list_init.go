package ui

import (
	"fmt"

	"charm.land/bubbles/v2/list"
	"github.com/charmbracelet/x/ansi"
	apiUtils "github.com/gohyuhan/rift/api/utils"
	"github.com/gohyuhan/rift/db"
	"github.com/gohyuhan/rift/i18n"
	pb "github.com/gohyuhan/rift/proto"
	"github.com/gohyuhan/rift/style"
	"go.etcd.io/bbolt"
	"google.golang.org/protobuf/proto"
)

// ----------------------------------
//
//	Reads every entry in the waypoint bucket and builds a display list.
//	Uses a read-only View transaction; any writes (corruption recording) are
//	deferred to a separate Update transaction after View completes.
//	Each waypoint occupies three consecutive lines in the returned slice:
//	  1. waypoint name (cyan for active; muted/faint + sealed label for sealed)
//	  2. waypoint path (blue-gray, faint, indented with two spaces)
//	  3. blank line separator
//	Corrupted proto data stops ForEach on the first affected entry; the name is
//	captured, recorded in the corrupted-records bucket via a follow-up Update,
//	and the caller receives a corruption-specific error for that waypoint.
//
// ---------------------------------
func getAllWaypointsInfo(bboltDb *bbolt.DB) ([]waypointInfo, error) {
	var waypointsInfo []waypointInfo
	var corruptedWaypointName []string
	waypointCorrupted := false

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

			// contruct the waypoint info type
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

	if waypointCorrupted {
		viewErr = apiUtils.RecordCorruptedWaypointInfo(bboltDb, corruptedWaypointName)
	}

	return waypointsInfo, viewErr
}

func initWaypointInfoListModel(m *WaypointInterctiveModel) error {
	previousSelectedWaypoint := m.WaypointInfoList.SelectedItem()
	selectedWayPointCursorPosition := -1

	latestWaypointInfoArray := []list.Item{}

	titleWidthLimit := m.Width - ListItemOrTitleWidthPad - 2

	allWaypointsInfo, err := getAllWaypointsInfo(m.BboltDb)
	if err != nil {
		return err
	}

	if previousSelectedWaypoint != nil {
		previousSelectedWaypointInfo := previousSelectedWaypoint.(waypointInfoItem)
		for index, waypoint := range allWaypointsInfo {
			if waypoint.WaypointName == previousSelectedWaypointInfo.WaypointName {
				selectedWayPointCursorPosition = index + 1
			}
			latestWaypointInfoArray = append(latestWaypointInfoArray, waypointInfoItem(waypoint))
		}
	} else {
		for _, waypoint := range allWaypointsInfo {
			latestWaypointInfoArray = append(latestWaypointInfoArray, waypointInfoItem(waypoint))
		}
	}

	m.WaypointInfoList = list.New(latestWaypointInfoArray, waypointInfoDelegate{}, m.Width, m.Height-5)
	m.WaypointInfoList.SetShowPagination(false)
	m.WaypointInfoList.SetShowStatusBar(false)
	m.WaypointInfoList.SetFilteringEnabled(false)
	m.WaypointInfoList.SetShowFilter(false)
	m.WaypointInfoList.SetShowHelp(true)

	m.WaypointInfoList.Title = ansi.Truncate(i18n.LANGUAGEMAPPING.WaypointInfoListTitle, titleWidthLimit, "...")
	m.WaypointInfoList.Styles.Title = style.NewStyle.Bold(true)
	m.WaypointInfoList.Styles.PaginationStyle = style.NewStyle
	m.WaypointInfoList.Styles.TitleBar = style.NewStyle
	m.WaypointInfoList.Styles.HelpStyle = style.NewStyle.MarginTop(0).MarginBottom(0).PaddingTop(0).PaddingBottom(0)

	if selectedWayPointCursorPosition >= 0 {
		m.WaypointInfoList.Select(selectedWayPointCursorPosition)
		m.WaypointInfoListCursorPosition = selectedWayPointCursorPosition
	} else {
		if m.WaypointInfoListCursorPosition > len(m.WaypointInfoList.Items())-1 {
			m.WaypointInfoList.Select(len(m.WaypointInfoList.Items()) - 1)
			m.WaypointInfoListCursorPosition = len(m.WaypointInfoList.Items()) - 1
		} else {
			m.WaypointInfoList.Select(m.WaypointInfoListCursorPosition)
		}
	}

	return nil
}
