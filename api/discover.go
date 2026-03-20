package api

import (
	"fmt"
	"time"

	"github.com/gohyuhan/rift/db"
	"github.com/gohyuhan/rift/i18n"
	"github.com/gohyuhan/rift/logger"
	pb "github.com/gohyuhan/rift/proto"
	"github.com/gohyuhan/rift/style"
	"github.com/gohyuhan/rift/utils"
	"github.com/spf13/cobra"
	"go.etcd.io/bbolt"
	"google.golang.org/protobuf/proto"
)

// ----------------------------------
//
//	Validates the current working directory, checks for reserved keywords,
//	opens the DB, and saves the given waypoint name mapped to the CWD.
//
// ----------------------------------
var RiftDiscoverFunc = func(command *cobra.Command, args []string) error {
	cwd, cwdErr := utils.GetCWD()
	if cwdErr != nil {
		return cwdErr
	}
	isDir, isDirErr := utils.CheckIsDir(cwd)
	if isDirErr != nil {
		return isDirErr
	}

	if !isDir {
		errorMessage := style.RenderStringWithColor(i18n.LANGUAGEMAPPING.CWDIsNotDirError, style.ColorError, false)
		return fmt.Errorf("%s", errorMessage)
	}

	if len(args) == 0 {
		return command.Help()
	}

	waypointName := args[0]

	if err := CheckIfKeywordIsReservedForRift(waypointName); err != nil {
		return err
	}

	bboltDB, bboltDBErr := db.OpenDB()
	if bboltDBErr != nil {
		return bboltDBErr
	}
	defer db.CloseDB(bboltDB)

	saveWaypointErr := saveWaypoint(bboltDB, waypointName, cwd)

	if saveWaypointErr != nil {
		return saveWaypointErr
	}

	successMessage := style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.RiftSavedWaypoint, waypointName, cwd), style.ColorGreenSoft, false)
	logger.LOGGER.LogToTerminal([]string{successMessage})

	return nil
}

// ----------------------------------
//
//	Persists a new waypoint entry into the DB bucket, rejecting duplicates.
//
// ----------------------------------
func saveWaypoint(bboltDb *bbolt.DB, waypointName string, path string) error {
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
			errorMessage := style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.RiftWaypointAlreadyExistsError, waypointName, existingWaypoint.WaypointPath), style.ColorError, false)
			return fmt.Errorf("%s", errorMessage)
		}

		waypoint := &pb.Waypoint{
			WaypointName:           waypointName,
			WaypointPath:           path,
			WaypointAddedAt:        time.Now().UTC().Format(time.RFC3339),
			WaypointTravelledCount: 0,
			WaypointIsSealed:       false,
		}

		data, err := proto.Marshal(waypoint)
		if err != nil {
			return err
		}

		return bucket.Put([]byte(waypointName), data)
	})
}
