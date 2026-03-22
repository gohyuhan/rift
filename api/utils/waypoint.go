// ----------------------------------
//
//	WAYPOINT RELATED UTILS
//
// ----------------------------------

package utils

import (
	"fmt"
	"strings"

	"github.com/gohyuhan/rift/db"
	"github.com/gohyuhan/rift/i18n"
	pb "github.com/gohyuhan/rift/proto"
	"github.com/gohyuhan/rift/style"
	"go.etcd.io/bbolt"
	"google.golang.org/protobuf/proto"
)

// ----------------------------------
//
//	Persists the waypoint name into the corrupted-records bucket via its own
//	Update transaction (so the write commits regardless of what the caller's
//	transaction did), then returns a user-facing corruption error.
//	The Update's own error is intentionally not propagated — recording is
//	best-effort; the corruption error is always returned to the caller.
//
// ----------------------------------
func RecordCorruptedWaypointInfo(bboltDB *bbolt.DB, corruptedWaypointsName []string) error {
	// best-effort write — ignore the Update error; the caller always gets the corruption message
	bboltDB.Update(func(tx *bbolt.Tx) error {
		waypointCorruptedBucket := tx.Bucket(db.WaypointDataCorruptedBucketRecord)
		if waypointCorruptedBucket != nil {
			for _, corruptedWaypoint := range corruptedWaypointsName {
				waypointCorruptedBucket.Put([]byte(corruptedWaypoint), []byte(corruptedWaypoint))
			}
		}
		return nil
	})
	return fmt.Errorf("%s", style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.WaypointDataCorruptedError, strings.Join(corruptedWaypointsName, ",")), style.ColorError, false))
}

// ----------------------------------
//
//	Fetches and deserializes the named waypoint from the bucket within an
//	already-open Update transaction. Returns the bucket, the deserialized
//	record, or an error if the bucket is missing, the waypoint does not exist,
//	or the stored proto is corrupted. Callers mutate the returned record and
//	re-persist it via bucket.Put.
//
// ----------------------------------
func GetWaypointForUpdate(tx *bbolt.Tx, waypointName string) (*bbolt.Bucket, *pb.Waypoint, error) {
	bucket := tx.Bucket(db.WaypointBucket)
	if bucket == nil {
		return nil, nil, fmt.Errorf("%s", style.RenderStringWithColor(i18n.LANGUAGEMAPPING.WaypointBucketNotFoundError, style.ColorError, false))
	}

	existing := bucket.Get([]byte(waypointName))
	if existing == nil {
		return nil, nil, fmt.Errorf("%s", style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.RiftWaypointDoNotExistsError, waypointName), style.ColorError, false))
	}

	waypoint := &pb.Waypoint{}
	if err := proto.Unmarshal(existing, waypoint); err != nil {
		return nil, nil, fmt.Errorf("%s", style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.WaypointDataCorruptedError, waypointName), style.ColorError, false))
	}

	return bucket, waypoint, nil
}

// ----------------------------------
//
//	Persists a mutated waypoint record back into its bucket. Returns an error
//	if marshalling or the bucket write fails.
//
// ----------------------------------
func PutWaypoint(bucket *bbolt.Bucket, waypointName string, waypoint *pb.Waypoint) error {
	data, err := proto.Marshal(waypoint)
	if err != nil {
		return fmt.Errorf("%s", style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.RiftWaypointUpdateError, waypointName), style.ColorError, false))
	}
	return bucket.Put([]byte(waypointName), data)
}

// ----------------------------------
//
//	Sets the sealed flag and reason on the named waypoint. A sealed waypoint is
//	one whose path no longer exists on disk; rift will refuse to travel to it
//	until it is explicitly unsealed. Called internally when a path-existence
//	check fails.
//
// ----------------------------------
func UpdateWaypointIsSeal(bboltDb *bbolt.DB, waypointName string, sealed bool, reason string) error {
	return bboltDb.Update(func(tx *bbolt.Tx) error {
		bucket, waypoint, err := GetWaypointForUpdate(tx, waypointName)
		if err != nil {
			return err
		}
		waypoint.WaypointIsSealed = sealed
		waypoint.WaypointSealedReason = reason
		return PutWaypoint(bucket, waypointName, waypoint)
	})
}

// ----------------------------------
//
//	Increments the travelled count for the named waypoint in the DB bucket.
//	Returns an error if the bucket is missing, the data is unreadable, or the
//	waypoint does not exist.
//
// ----------------------------------
func UpdateWaypointTravelledCount(bboltDb *bbolt.DB, waypointName string) error {
	return bboltDb.Update(func(tx *bbolt.Tx) error {
		bucket, waypoint, err := GetWaypointForUpdate(tx, waypointName)
		if err != nil {
			return err
		}
		waypoint.WaypointTravelledCount += 1
		return PutWaypoint(bucket, waypointName, waypoint)
	})
}
