package api

import (
	"fmt"
	"strings"
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
//	Cobra handler for the discover command.
//	Resolves the current working directory, validates it is a real directory,
//	guards against reserved keyword names, then persists the waypoint mapping
//	(name → CWD) into the DB.
//
// ----------------------------------
var RiftDiscoverFunc = func(command *cobra.Command, args []string) error {
	// resolve and validate the current working directory
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

	// no waypoint name provided — show help
	if len(args) == 0 {
		return command.Help()
	}

	waypointName := strings.TrimSpace(args[0])

	// reject names that clash with rift's own subcommands
	if err := CheckIfKeywordIsReservedForRift(waypointName); err != nil {
		return err
	}

	// open DB and persist the new waypoint
	bboltDB, bboltDBErr := db.OpenDB()
	if bboltDBErr != nil {
		return bboltDBErr
	}
	defer db.CloseDB(bboltDB)

	saveWaypointErr := saveWaypoint(bboltDB, waypointName, cwd)

	if saveWaypointErr != nil {
		return saveWaypointErr
	}

	message := style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.RiftSavedWaypoint, waypointName, cwd), style.ColorGreenSoft, false)
	logger.LOGGER.LogToTerminal([]string{message})

	return nil
}

// ----------------------------------
//
//	Persists a new waypoint (name → path) into the waypoint bucket.
//	Rejects the write if a waypoint with the same name already exists,
//	whether healthy or corrupted.
//
// ----------------------------------
func saveWaypoint(bboltDb *bbolt.DB, waypointName string, path string) error {
	return bboltDb.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(db.WaypointBucket)
		if bucket == nil {
			return fmt.Errorf("%s", style.RenderStringWithColor(i18n.LANGUAGEMAPPING.WaypointBucketNotFoundError, style.ColorError, false))
		}

		// duplicate guard: reject if a waypoint with this name already exists
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

		// build the new waypoint record with defaults
		waypoint := &pb.Waypoint{
			WaypointName:           waypointName,
			WaypointPath:           path,
			WaypointAddedAt:        time.Now().UTC().Format(time.RFC3339),
			WaypointTravelledCount: 0,
			WaypointIsSealed:       false,
			WaypointSealedReason:   "",
		}

		data, err := proto.Marshal(waypoint)
		if err != nil {
			return err
		}

		return bucket.Put([]byte(waypointName), data)
	})
}
