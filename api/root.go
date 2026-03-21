package api

import (
	"fmt"
	"strings"

	"github.com/gohyuhan/rift/db"
	"github.com/gohyuhan/rift/i18n"
	pb "github.com/gohyuhan/rift/proto"
	"github.com/gohyuhan/rift/settings"
	"github.com/gohyuhan/rift/style"
	"github.com/gohyuhan/rift/updater"
	"github.com/gohyuhan/rift/utils"
	"github.com/spf13/cobra"
	"go.etcd.io/bbolt"
	"google.golang.org/protobuf/proto"
)

// ----------------------------------
//
//	Cobra handler for the root rift command.
//	If --update is passed, triggers an immediate update check and returns.
//	When called with no args, it handles settings flags (language, autoupdate,
//	download-pre-release); if no flags were passed it falls through to help.
//	When called with an arg, it travels to the named waypoint by printing a
//	cd command for the shell wrapper to eval, then increments the travel count.
//
// ----------------------------------
var RiftRootFunc = func(cmd *cobra.Command, args []string) error {
	if cmd.Flags().Changed("update") {
		updater.Update()
		return nil
	}

	if len(args) < 1 {
		// settings flags can only take effect when there are no waypoint args
		languageFlagCalled := cmd.Flags().Changed("language")
		autoupdateFlagCalled := cmd.Flags().Changed("autoupdate")
		downloadPreReleaseFlagCalled := cmd.Flags().Changed("download-pre-release")

		if languageFlagCalled {
			languageSetting, languageSettingErr := cmd.Flags().GetString("language")
			if languageSettingErr != nil {
				return fmt.Errorf("%s", style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.RiftFlagRetrieveError, "language", languageSettingErr.Error()), style.ColorError, false))
			}
			settings.UpdateLanguageCode(languageSetting)
		}

		if autoupdateFlagCalled {
			autoUpdateSetting, autoUpdateSettingErr := cmd.Flags().GetBool("autoupdate")
			if autoUpdateSettingErr != nil {
				return fmt.Errorf("%s", style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.RiftFlagRetrieveError, "autoupdate", autoUpdateSettingErr.Error()), style.ColorError, false))
			}
			settings.UpdateAutoUpdate(autoUpdateSetting)
		}

		if downloadPreReleaseFlagCalled {
			downloadPreReleaseSetting, downloadPreReleaseSettingErr := cmd.Flags().GetBool("download-pre-release")
			if downloadPreReleaseSettingErr != nil {
				return fmt.Errorf("%s", style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.RiftFlagRetrieveError, "download-pre-release", downloadPreReleaseSettingErr.Error()), style.ColorError, false))
			}
			settings.UpdateDownloadPreRelease(downloadPreReleaseSetting)
		}

		// no flags were passed — show help
		if !(languageFlagCalled || autoupdateFlagCalled || downloadPreReleaseFlagCalled) {
			return cmd.Help()
		}

		return nil
	}

	waypointName := strings.TrimSpace(args[0])

	// open DB for reading waypoint data
	bboltDB, bboltDBErr := db.OpenDB()
	if bboltDBErr != nil {
		return bboltDBErr
	}
	defer db.CloseDB(bboltDB)

	// look up the waypoint path; errors here are user-visible (sealed, missing, corrupted)
	retrievedPath, retrieveErr := retrieveWaypointInfo(bboltDB, waypointName)
	if retrieveErr != nil {
		return retrieveErr
	}

	// Only this line goes to stdout — the shell wrapper evals it.
	fmt.Printf("cd %q", retrievedPath)

	// best-effort: increment travel count; failure is silently ignored
	updateWaypointTravelledCount(bboltDB, waypointName)

	return nil
}

// ----------------------------------
//
//	Looks up a waypoint by name and validates it is travelable.
//	Uses a read-only View transaction for the lookup; sealing on a missing path
//	is deferred to a second write transaction after View completes — bbolt write
//	locks are not reentrant, so calling Update inside View (or Update) would deadlock.
//	Returns the stored path, or an error when:
//	  - the waypoint bucket is missing
//	  - the waypoint does not exist
//	  - the stored proto data is corrupted
//	  - the waypoint is already sealed
//	  - the waypoint path no longer exists on disk (seals it via a follow-up write tx)
//
// ----------------------------------
func retrieveWaypointInfo(bboltDb *bbolt.DB, waypointName string) (string, error) {
	retrievedPath := ""
	waypointCorrupted := false
	needToSealWaypoint := false
	needToSealReason := ""

	viewErr := bboltDb.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(db.WaypointBucket)
		if bucket == nil {
			return fmt.Errorf("%s", style.RenderStringWithColor(i18n.LANGUAGEMAPPING.WaypointBucketNotFoundError, style.ColorError, false))
		}

		// check the waypoint exists in the bucket
		existing := bucket.Get([]byte(waypointName))
		if existing == nil {
			errorMessage := style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.RiftWaypointDoNotExistsError, waypointName), style.ColorError, false)
			return fmt.Errorf("%s", errorMessage)
		}

		// deserialize the stored proto; set the flag and return nil so the View
		// commits cleanly — corruption recording is deferred to a follow-up Update
		existingWaypoint := &pb.Waypoint{}
		protoErr := proto.Unmarshal(existing, existingWaypoint)
		if protoErr != nil {
			waypointCorrupted = true
			return nil
		}

		// sealed means the path no longer exists or was manually sealed; block travel
		if existingWaypoint.WaypointIsSealed {
			return fmt.Errorf("%s", style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.RiftWaypointSealedError, waypointName, existingWaypoint.WaypointSealedReason), style.ColorError, false))
		}

		// verify the path still exists on disk; if not, seal the waypoint and abort
		isPathExist, isPathExistErr := utils.CheckIsPathExist(existingWaypoint.WaypointPath)
		if !isPathExist {
			needToSealWaypoint = true
			needToSealReason = isPathExistErr.Error()
			return fmt.Errorf("%s", style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.RiftWaypointSealedError, waypointName, isPathExistErr.Error()), style.ColorError, false))
		}

		retrievedPath = existingWaypoint.WaypointPath
		return nil
	})

	if waypointCorrupted {
		viewErr = recordCorruptedWaypointInfo(bboltDb, waypointName)
	}

	if needToSealWaypoint {
		updateWaypointIsSeal(bboltDb, waypointName, true, needToSealReason)
	}

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
		bucket, waypoint, err := getWaypointForUpdate(tx, waypointName)
		if err != nil {
			return err
		}
		waypoint.WaypointTravelledCount += 1
		return putWaypoint(bucket, waypointName, waypoint)
	})
}

// ----------------------------------
//
//	Sets the sealed flag and reason on the named waypoint. A sealed waypoint is
//	one whose path no longer exists on disk; rift will refuse to travel to it
//	until it is explicitly unsealed. Called internally when a path-existence
//	check fails.
//
// ----------------------------------
func updateWaypointIsSeal(bboltDb *bbolt.DB, waypointName string, sealed bool, reason string) error {
	return bboltDb.Update(func(tx *bbolt.Tx) error {
		bucket, waypoint, err := getWaypointForUpdate(tx, waypointName)
		if err != nil {
			return err
		}
		waypoint.WaypointIsSealed = sealed
		waypoint.WaypointSealedReason = reason
		return putWaypoint(bucket, waypointName, waypoint)
	})
}
