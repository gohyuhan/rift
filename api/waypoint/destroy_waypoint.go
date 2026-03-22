package waypoint

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
//	Permanently removes the named waypoint from the waypoint bucket.
//	Uses a write Update transaction. The operation is idempotent: if the
//	waypoint does not exist, it is treated as already destroyed (success).
//	Fails only when the bucket itself is missing or bbolt returns a hard error.
//
// ----------------------------------
func destroyDiscoveredWaypoint(bboltDb *bbolt.DB, waypointName string) error {
	return bboltDb.Update(func(tx *bbolt.Tx) error {
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

		// report the destruction to the terminal
		message := style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.RiftWaypointDestroySuccess, waypointName), style.ColorGreenSoft, false)
		logger.LOGGER.LogToTerminal([]string{message})

		return nil
	})
}
