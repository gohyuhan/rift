package discover

import (
	"fmt"
	"time"

	"github.com/gohyuhan/rift/db"
	"github.com/gohyuhan/rift/i18n"
	pb "github.com/gohyuhan/rift/proto"
	"github.com/gohyuhan/rift/style"
	"go.etcd.io/bbolt"
	"google.golang.org/protobuf/proto"
)

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
