package features

import (
	"fmt"

	"github.com/gohyuhan/rift/db"
	"github.com/gohyuhan/rift/i18n"
	"github.com/gohyuhan/rift/logger"
	"github.com/gohyuhan/rift/style"
	"go.etcd.io/bbolt"
)

// ----------------------------------
//
//	Permanently removes the named waypoint from the waypoint bucket and its
//	corrupted-data record if one exists, within a single write transaction.
//	The operation is idempotent: if the waypoint does not exist it is treated
//	as already destroyed (success). Fails only when the waypoint bucket itself
//	is missing or bbolt returns a hard error.
//
// ----------------------------------
func DestroyDiscoveredWaypoint(waypointName string, logToTerminal bool) error {
	bboltWriteDb, bboltWriteDbErr := db.OpenWriteDB()
	if bboltWriteDbErr != nil {
		return bboltWriteDbErr
	}
	defer db.CloseDB(bboltWriteDb)

	dbErr := bboltWriteDb.Update(func(tx *bbolt.Tx) error {
		// ensure the waypoint bucket exists before attempting the delete
		waypointBucket := tx.Bucket(db.WaypointBucket)
		if waypointBucket == nil {
			return fmt.Errorf("%s", style.RenderStringWithColor(i18n.LANGUAGEMAPPING.WaypointBucketNotFoundError, style.ColorError, false))
		}

		// bbolt.Delete returns nil for missing keys; the behavior here is
		// intentionally idempotent — destroying a non-existent waypoint succeeds
		destroyWaypointErr := waypointBucket.Delete([]byte(waypointName))
		if destroyWaypointErr != nil {
			return fmt.Errorf("%s", style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.RiftWaypointDestroyError, waypointName, destroyWaypointErr.Error()), style.ColorError, false))
		}

		// also remove the corrupted-data record for this waypoint if one exists
		corruptedWaypointBucket := tx.Bucket(db.WaypointDataCorruptedBucketRecord)
		if corruptedWaypointBucket != nil {
			corruptedWaypointBucket.Delete([]byte(waypointName))
		}

		return nil
	})

	// report the destruction to the terminal
	if dbErr == nil && logToTerminal {
		message := style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.RiftWaypointDestroySuccess, waypointName), style.ColorGreenSoft, false)
		logger.LOGGER.LogToTerminal([]string{message})
	}

	return dbErr
}
