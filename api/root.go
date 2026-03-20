package api

import (
	"fmt"

	"github.com/gohyuhan/rift/db"
	"github.com/gohyuhan/rift/i18n"
	pb "github.com/gohyuhan/rift/proto"
	"github.com/gohyuhan/rift/style"
	"github.com/spf13/cobra"
	"go.etcd.io/bbolt"
	"google.golang.org/protobuf/proto"
)

// ----------------------------------
//
//	Cobra handler for the root rift command. Currently a placeholder.
//
// ----------------------------------
var RiftRootFunc = func(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return cmd.Help()
	}

	waypointName := args[0]

	bboltDB, bboltDBErr := db.OpenDB()
	if bboltDBErr != nil {
		return bboltDBErr
	}
	defer db.CloseDB(bboltDB)

	retrievedPath, retrieveErr := retrieveWaypointInfo(bboltDB, waypointName)
	if retrieveErr != nil {
		return retrieveErr
	}

	// Only this line goes to stdout — the shell wrapper evals it.
	fmt.Printf("cd %q", retrievedPath)

	updateWaypointTravelledCount(bboltDB, waypointName)

	return nil
}

// ----------------------------------
//
//	Looks up a waypoint by name in the DB bucket. Returns the waypoint path
//	or an error if the bucket is missing, the data is unreadable, or the
//	waypoint does not exist.
//
// ----------------------------------
func retrieveWaypointInfo(bboltDb *bbolt.DB, waypointName string) (string, error) {
	retrievedPath := ""
	viewErr := bboltDb.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(db.WaypointBucket)
		if bucket == nil {
			return fmt.Errorf("%s", style.RenderStringWithColor(i18n.LANGUAGEMAPPING.WaypointBucketNotFoundError, style.ColorError, false))
		}

		existing := bucket.Get([]byte(waypointName))
		if existing == nil {
			errorMessage := style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.RiftWaypointDoNotExistsError, waypointName), style.ColorError, false)
			return fmt.Errorf("%s", errorMessage)
		}

		existingWaypoint := &pb.Waypoint{}
		protoErr := proto.Unmarshal(existing, existingWaypoint)
		if protoErr != nil {
			waypointCorruptedBucket := tx.Bucket(db.WaypointDataCorruptedBucketRecord)
			if waypointCorruptedBucket != nil {
				waypointCorruptedBucket.Put([]byte(waypointName), []byte(waypointName))
			}
			return fmt.Errorf("%s", style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.WaypointDataCorruptedError, waypointName), style.ColorError, false))
		}

		retrievedPath = existingWaypoint.WaypointPath
		return nil
	})

	return retrievedPath, viewErr
}

// ----------------------------------
//
//	Increments the travelled count for the named waypoint in the DB bucket.
//	Returns an error if the bucket is missing, the data is unreadable, or the
//	waypoint does not exist.
//
// ----------------------------------
func updateWaypointTravelledCount(bboltDb *bbolt.DB, waypointName string) error {
	return bboltDb.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(db.WaypointBucket)
		if bucket == nil {
			return fmt.Errorf("%s", style.RenderStringWithColor(i18n.LANGUAGEMAPPING.WaypointBucketNotFoundError, style.ColorError, false))
		}

		existing := bucket.Get([]byte(waypointName))
		if existing != nil {
			existingWaypoint := &pb.Waypoint{}
			protoErr := proto.Unmarshal(existing, existingWaypoint)
			if protoErr != nil {
				return fmt.Errorf("%s", style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.WaypointDataCorruptedError, waypointName), style.ColorError, false))
			}

			existingWaypoint.WaypointTravelledCount += 1
			updatedWaypointInfo, updateErr := proto.Marshal(existingWaypoint)

			if updateErr != nil {
				errorMessage := style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.RiftWaypointUpdateError, waypointName), style.ColorError, false)
				return fmt.Errorf("%s", errorMessage)
			}

			return bucket.Put([]byte(waypointName), updatedWaypointInfo)
		}
		errorMessage := style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.RiftWaypointDoNotExistsError, waypointName), style.ColorError, false)
		return fmt.Errorf("%s", errorMessage)
	})
}
