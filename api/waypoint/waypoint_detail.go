package waypoint

import (
	"fmt"
	"strconv"
	"time"

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
//	Fetches a single named waypoint and builds a detail display list.
//	Uses a read-only View transaction; corruption recording is deferred to a
//	separate Update transaction after View completes (same pattern as
//	retrieveAllWaypointInfo).
//	The returned slice has one labelled row per field; the sealed-reason row is
//	only included when WaypointIsSealed is true.
//	Labels are padded to uniform width via style.PadAndRenderLabels so values
//	align in a clean column regardless of the active language.
//
// ----------------------------------
func retrieveWaypointInfoDetail(bboltDb *bbolt.DB, waypointName string) ([]string, error) {
	var waypointDetailInfo []string
	waypointCorrupted := false

	viewErr := bboltDb.View(func(tx *bbolt.Tx) error {
		// ensure the waypoint bucket exists before looking up the key
		waypointBucket := tx.Bucket(db.WaypointBucket)
		if waypointBucket == nil {
			return fmt.Errorf("%s", style.RenderStringWithColor(i18n.LANGUAGEMAPPING.WaypointBucketNotFoundError, style.ColorError, false))
		}

		// look up the specific waypoint by name
		existing := waypointBucket.Get([]byte(waypointName))
		if existing == nil {
			errorMessage := style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.RiftWaypointDoNotExistsError, waypointName), style.ColorError, false)
			return fmt.Errorf("%s", errorMessage)
		}

		// deserialize the stored proto; flag corruption and exit the View cleanly —
		// the actual recording write is deferred to after the View completes
		existingWaypoint := &pb.Waypoint{}
		protoErr := proto.Unmarshal(existing, existingWaypoint)
		if protoErr != nil {
			waypointCorrupted = true
			return nil
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

		return nil
	})

	// View is complete — safe to open an Update for the corruption write
	if waypointCorrupted {
		viewErr = apiUtils.RecordCorruptedWaypointInfo(bboltDb, waypointName)
	}

	return waypointDetailInfo, viewErr
}
