package features

import (
	"fmt"
	"time"

	apiUtils "github.com/gohyuhan/rift/api/utils"
	"github.com/gohyuhan/rift/db"
	"github.com/gohyuhan/rift/i18n"
	"github.com/gohyuhan/rift/logger"
	"github.com/gohyuhan/rift/style"
	"github.com/gohyuhan/rift/utils"
	"go.etcd.io/bbolt"
)

// ----------------------------------
//
//	Reassigns the named waypoint to a new path and resets its state.
//	If rebindTo is empty, defaults to the current working directory.
//	Validates that rebindTo is an existing directory before writing.
//	On success, clears the sealed flag, sealed reason, and travelled count,
//	and updates the discovered timestamp to now (UTC).
//
// ----------------------------------
func RebindWaypoint(waypointName string, rebindTo string, logToTerminal bool) error {
	// validate if the rebindTo is valid or not before opening the Update transaction;
	// this avoids unnecessary DB writes if the path is invalid
	// and also prevent holding the DB lock during potentially slow filesystem operations
	if rebindTo != "" {
		isDirExist, isDirExistErr := utils.CheckIsDir(rebindTo)
		if isDirExistErr != nil {
			return isDirExistErr
		} else if !isDirExist {
			return fmt.Errorf("%s", style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.RiftWaypointRebindNotDirError, rebindTo), style.ColorError, false))
		}
	} else {
		cwd, getCWDErr := utils.GetCWD()
		if getCWDErr != nil {
			return getCWDErr
		}
		rebindTo = cwd
	}

	bboltWriteDb, bboltWriteDbErr := db.OpenWriteDB()
	if bboltWriteDbErr != nil {
		return bboltWriteDbErr
	}
	defer db.CloseDB(bboltWriteDb)

	return bboltWriteDb.Update(func(tx *bbolt.Tx) error {
		// fetch the current waypoint record and its bucket in a single helper call
		waypointBucket, waypoint, retrieveErr := apiUtils.GetWaypointForUpdate(tx, waypointName)
		if retrieveErr != nil {
			return retrieveErr
		}

		// update all mutable fields, clear rune; clear sealed state so the waypoint is immediately usable
		waypoint.WaypointPath = rebindTo
		waypoint.WaypointIsSealed = false
		waypoint.WaypointSealedReason = ""
		waypoint.WaypointTravelledCount = 0
		waypoint.WaypointAddedAt = time.Now().UTC().Format(time.RFC3339)
		waypoint.EnterRunes = nil
		waypoint.LeaveRunes = nil

		// persist the updated record back under the same key
		putWaypointErr := apiUtils.PutWaypoint(waypointBucket, waypointName, waypoint)
		if putWaypointErr != nil {
			return putWaypointErr
		}

		// report the new binding to the terminal
		if logToTerminal {
			message := style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.RiftWaypointRebindSuccess, waypointName, rebindTo), style.ColorGreenSoft, false)
			logger.LOGGER.LogToTerminal([]string{message})
		}
		return nil
	})
}
