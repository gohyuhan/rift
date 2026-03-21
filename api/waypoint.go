package api

import (
	"fmt"

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

	// named waypoint detail view (future)
	return nil
}

// ----------------------------------
//
//	Reads every entry in the waypoint bucket and builds a display list.
//	Each waypoint occupies three consecutive lines in the returned slice:
//	  1. waypoint name (cyan for active; muted/faint + sealed label for sealed)
//	  2. waypoint path (blue-gray, faint, indented with two spaces)
//	  3. blank line separator
//	Corrupted proto data is recorded in the corrupted-records bucket and
//	causes ForEach to short-circuit with an error; the caller then receives
//	a generic retrieve-all error.
//
// ----------------------------------
func retrieveAllWaypointInfo(bboltDb *bbolt.DB) ([]string, error) {
	var waypointsInfo []string
	dbError := bboltDb.View(func(tx *bbolt.Tx) error {
		// ensure the waypoint bucket exists before iterating
		waypointBucket := tx.Bucket(db.WaypointBucket)
		if waypointBucket == nil {
			return fmt.Errorf("%s", style.RenderStringWithColor(i18n.LANGUAGEMAPPING.WaypointBucketNotFoundError, style.ColorError, false))
		}

		// walk every key-value pair in the bucket
		retrieveError := waypointBucket.ForEach(func(k, v []byte) error {
			// deserialize the stored proto; record and abort on corruption
			existingWaypoint := &pb.Waypoint{}
			protoErr := proto.Unmarshal(v, existingWaypoint)
			if protoErr != nil {
				return recordCorruptedWaypointInfo(tx, string(k))
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
	return waypointsInfo, dbError
}
