package features

import (
	"strconv"
	"time"

	apiUtils "github.com/gohyuhan/rift/api/utils"
	"github.com/gohyuhan/rift/db"
	"github.com/gohyuhan/rift/i18n"
	"github.com/gohyuhan/rift/style"
)

// ----------------------------------
//
//	Fetches a single named waypoint via apiUtils.RetrieveWaypointInfo and builds
//	a detail display list. The returned slice has one labelled row per field:
//	name, path, discovered-at, travelled count, sealed state, and (when sealed)
//	sealed reason. The name value includes two \uf4bf icons — the first for the
//	enter-rune slot and the second for the leave-rune slot — rendered in vibrant
//	cyan when the slot is populated or muted/faint when it is empty. Labels are
//	padded to uniform width via style.PadAndRenderLabels so values align in a
//	clean column regardless of the active language.
//
// ----------------------------------
func RetrieveWaypointInfoDetail(waypointName string) ([]string, error) {
	var waypointDetailInfo []string
	var viewErr error
	waypointCorrupted := false

	// open DB so we can read waypoint records
	bboltReadDb, bboltReadDbErr := db.OpenReadDB()
	if bboltReadDbErr != nil {
		return waypointDetailInfo, bboltReadDbErr
	}
	defer db.CloseDB(bboltReadDb)

	existingWaypoint, viewErr := apiUtils.RetrieveWaypointInfo(waypointName)
	if viewErr != nil {
		return waypointDetailInfo, viewErr
	}

	// collect raw labels and pad them to uniform width
	rawLabels := []string{
		i18n.LANGUAGEMAPPING.RiftWaypointDetailName,
		i18n.LANGUAGEMAPPING.RiftWaypointDetailPath,
		i18n.LANGUAGEMAPPING.RiftWaypointDetailDiscovered,
		i18n.LANGUAGEMAPPING.RiftWaypointDetailTravelledCount,
		i18n.LANGUAGEMAPPING.RiftWaypointDetailSealed,
	}
	if existingWaypoint.WaypointIsSealed {
		rawLabels = append(rawLabels, i18n.LANGUAGEMAPPING.RiftWaypointDetailSealedReason)
	}
	paddedLabels := style.PadAndRenderLabels(rawLabels, style.ColorBlueGrayMuted, true)

	// --- waypoint name ---
	var nameValue string
	if existingWaypoint.WaypointIsSealed {
		nameValue = style.RenderStringWithColor(waypointName, style.ColorSealedMuted, false)
	} else {
		nameValue = style.RenderStringWithColor(waypointName, style.ColorCyanSoft, false)
	}

	if existingWaypoint.EnterRunes != nil || len(existingWaypoint.EnterRunes) > 0 {
		nameValue = nameValue + " " + style.RenderStringWithColor("\uf4bf", style.ColorCyanSoft, false)
	} else {
		nameValue = nameValue + " " + style.RenderStringWithColor("\uf4bf", style.ColorSealedMuted, true)
	}

	if existingWaypoint.LeaveRunes != nil || len(existingWaypoint.LeaveRunes) > 0 {
		nameValue = nameValue + " " + style.RenderStringWithColor("\uf4bf", style.ColorCyanSoft, false)
	} else {
		nameValue = nameValue + " " + style.RenderStringWithColor("\uf4bf", style.ColorSealedMuted, true)
	}

	waypointDetailInfo = append(waypointDetailInfo, paddedLabels[0]+"  "+nameValue)

	// --- waypoint path ---
	pathValue := style.RenderStringWithColor(existingWaypoint.WaypointPath, style.ColorBlueMuted, false)
	waypointDetailInfo = append(waypointDetailInfo, paddedLabels[1]+"  "+pathValue)

	// --- waypoint discovered (UTC → local) ---
	discoveredDisplay := existingWaypoint.WaypointAddedAt
	if parsed, err := time.Parse(time.RFC3339, existingWaypoint.WaypointAddedAt); err == nil {
		discoveredDisplay = parsed.Local().Format("2006-01-02 15:04:05")
	}
	discoveredValue := style.RenderStringWithColor(discoveredDisplay, style.ColorBlueMuted, false)
	waypointDetailInfo = append(waypointDetailInfo, paddedLabels[2]+"  "+discoveredValue)

	// --- waypoint travelled count ---
	travelledValue := style.RenderStringWithColor(strconv.FormatInt(existingWaypoint.WaypointTravelledCount, 10), style.ColorBlueMuted, false)
	waypointDetailInfo = append(waypointDetailInfo, paddedLabels[3]+"  "+travelledValue)

	// --- waypoint sealed ---
	var sealedValue string
	if existingWaypoint.WaypointIsSealed {
		sealedValue = style.RenderStringWithColor(i18n.LANGUAGEMAPPING.RiftWaypointDetailSealedTrue, style.ColorSealedMuted, false)
	} else {
		sealedValue = style.RenderStringWithColor(i18n.LANGUAGEMAPPING.RiftWaypointDetailSealedFalse, style.ColorBlueMuted, false)
	}
	waypointDetailInfo = append(waypointDetailInfo, paddedLabels[4]+"  "+sealedValue)

	// --- sealed reason (only when sealed) ---
	if existingWaypoint.WaypointIsSealed {
		reasonValue := style.RenderStringWithColor(existingWaypoint.WaypointSealedReason, style.ColorSealedMuted, false)
		waypointDetailInfo = append(waypointDetailInfo, paddedLabels[5]+"  "+reasonValue)
	}

	// close early so it will not block write connection below
	db.CloseDB(bboltReadDb)

	// View is complete — safe to open an Update for the corruption write
	if waypointCorrupted {
		viewErr = apiUtils.RecordCorruptedWaypointInfo([]string{waypointName})
	}

	return waypointDetailInfo, viewErr
}
