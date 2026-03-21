package api

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gohyuhan/rift/db"
	"github.com/gohyuhan/rift/i18n"
	"github.com/gohyuhan/rift/logger"
	pb "github.com/gohyuhan/rift/proto"
	"github.com/gohyuhan/rift/style"
	"github.com/spf13/cobra"
	"go.etcd.io/bbolt"
	"google.golang.org/protobuf/proto"
)

// ----------------------------------
//
//	Cobra handler for the waypoint command.
//	With no args, lists every stored waypoint (name + path, sealed state).
//	With an arg, shows detailed info for the named waypoint (not yet implemented).
//
// ----------------------------------
var RiftWaypointFunc = func(cmd *cobra.Command, args []string) error {
	// open DB so we can read waypoint records
	bboltDB, bboltDBErr := db.OpenDB()
	if bboltDBErr != nil {
		return bboltDBErr
	}
	defer db.CloseDB(bboltDB)

	// no args — list all waypoints
	if len(args) < 1 {
		allWaypointInfo, allWaypointInfoErr := retrieveAllWaypointInfo(bboltDB)
		if allWaypointInfoErr != nil {
			return allWaypointInfoErr
		}

		logger.LOGGER.LogToTerminal(allWaypointInfo)

		return nil
	}

	waypointName := args[0]

	destroyFlagCalled := cmd.Flags().Changed("destroy")
	rebindFlagCalled := cmd.Flags().Changed("rebind")
	reforgeFlagCalled := cmd.Flags().Changed("reforge")

	if destroyFlagCalled {
		// if destroy flag is called, we destroy the discovered waypoint in the waypoint bucket
		destroyWaypointErr := destroyDiscoveredWaypoint(bboltDB, waypointName)
		if destroyWaypointErr != nil {
			return destroyWaypointErr
		}
	} else if rebindFlagCalled {
	} else if reforgeFlagCalled {
	} else {
		retrieveWaypointInfoDetail, retrieveWaypointInfoDetailErr := retrieveWaypointInfoDetail(bboltDB, waypointName)
		if retrieveWaypointInfoDetailErr != nil {
			return retrieveWaypointInfoDetailErr
		}

		logger.LOGGER.LogToTerminal(retrieveWaypointInfoDetail)
	}

	return nil
}

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
// ----------------------------------
func retrieveAllWaypointInfo(bboltDb *bbolt.DB) ([]string, error) {
	var waypointsInfo []string
	var corruptedWaypointName string
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
			if protoErr != nil {
				waypointCorrupted = true
				corruptedWaypointName = string(k)
				return fmt.Errorf("")
			}

			// build the waypoint name line; sealed entries use a dark dormant palette
			var waypointName string
			if existingWaypoint.WaypointIsSealed {
				// faint name + non-faint sealed label — both in the muted sealed color
				waypointName = style.RenderStringWithColor(string(k), style.ColorSealedMuted, true) +
					" " + style.RenderStringWithColor(i18n.LANGUAGEMAPPING.RiftWaypointSealedLabel, style.ColorSealedMuted, false)
			} else {
				waypointName = style.RenderStringWithColor(string(k), style.ColorCyanSoft, false)
			}

			// append name, indented path, then a blank separator line
			waypointsInfo = append(waypointsInfo, waypointName)
			waypointsInfo = append(waypointsInfo, style.RenderStringWithColor("  "+existingWaypoint.WaypointPath, style.ColorBlueGrayMuted, true))
			waypointsInfo = append(waypointsInfo, "")
			return nil
		})

		// wrap ForEach failure into a user-facing message
		if retrieveError != nil {
			return fmt.Errorf("%s", style.RenderStringWithColor(i18n.LANGUAGEMAPPING.RiftWaypointRetrieveAllError, style.ColorError, false))
		}

		return nil
	})

	if waypointCorrupted {
		viewErr = recordCorruptedWaypointInfo(bboltDb, corruptedWaypointName)
	}

	return waypointsInfo, viewErr
}

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
		viewErr = recordCorruptedWaypointInfo(bboltDb, waypointName)
	}

	return waypointDetailInfo, viewErr
}

func destroyDiscoveredWaypoint(bboltDb *bbolt.DB, waypointName string) error {
	return bboltDb.Update(func(tx *bbolt.Tx) error {
		// ensure the waypoint bucket exists before looking up the key
		waypointBucket := tx.Bucket(db.WaypointBucket)
		if waypointBucket == nil {
			return fmt.Errorf("%s", style.RenderStringWithColor(i18n.LANGUAGEMAPPING.WaypointBucketNotFoundError, style.ColorError, false))
		}

		destroyWaypointErr := waypointBucket.Delete([]byte(waypointName))

		if destroyWaypointErr != nil {
			return fmt.Errorf("%s", style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.RiftWaypointDestroyError, waypointName, destroyWaypointErr.Error()), style.ColorError, false))
		}
		message := style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.RiftWaypointDestroySuccess, waypointName), style.ColorGreenSoft, false)
		logger.LOGGER.LogToTerminal([]string{message})

		return nil
	})
}
