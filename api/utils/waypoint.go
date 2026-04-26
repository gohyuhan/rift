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
	"github.com/gohyuhan/rift/utils"
	"go.etcd.io/bbolt"
	"google.golang.org/protobuf/proto"
)

// ----------------------------------
//
//	raw data record for a single waypoint as read from the database;
//	mirrors the proto fields relevant to the list UI
//
// ----------------------------------
type WaypointInfo struct {
	WaypointName           string
	WaypointPath           string
	WaypointAddedAt        string
	WaypointTravelledCount int64
	WaypointIsSealed       bool
	WaypointSealedReason   string
	EnterRune              []*pb.RuneCmds
	LeaveRune              []*pb.RuneCmds
}

// ----------------------------------
//
//	Persists the waypoint name into the corrupted-records bucket via its own
//	Update transaction (so the write commits regardless of what the caller's
//	transaction did), then returns a user-facing corruption error.
//	The Update's own error is intentionally not propagated — recording is
//	best-effort; the corruption error is always returned to the caller.
//
// ----------------------------------
func RecordCorruptedWaypointInfo(corruptedWaypointsName []string) error {
	// best-effort write — ignore the Update error; the caller always gets the corruption message
	bboltWriteDb, bboltWriteDbErr := db.OpenWriteDB()
	if bboltWriteDbErr != nil {
		return bboltWriteDbErr
	}
	defer db.CloseDB(bboltWriteDb)

	bboltWriteDb.Update(func(tx *bbolt.Tx) error {
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
		return fmt.Errorf("%s", style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.RiftWaypointUpdateError, waypointName, err.Error()), style.ColorError, false))
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
func UpdateWaypointIsSeal(waypointName string, sealed bool, reason string) error {
	bboltWriteDb, bboltWriteDbErr := db.OpenWriteDB()
	if bboltWriteDbErr != nil {
		return bboltWriteDbErr
	}
	defer db.CloseDB(bboltWriteDb)

	return bboltWriteDb.Update(func(tx *bbolt.Tx) error {
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
//	Clears the sealed flag and sealed reason on the named waypoint.
//	After unsealing, rift will reattempt path validation on the next
//	navigation; if the path is still missing it will be re-sealed immediately.
//	Returns an error if the bucket is missing, the data is unreadable,
//	or the waypoint does not exist.
//
// ----------------------------------
func UpdateWaypointUnSeal(waypointName string) error {
	bboltWriteDb, bboltWriteDbErr := db.OpenWriteDB()
	if bboltWriteDbErr != nil {
		return bboltWriteDbErr
	}
	defer db.CloseDB(bboltWriteDb)

	return bboltWriteDb.Update(func(tx *bbolt.Tx) error {
		bucket, waypoint, err := GetWaypointForUpdate(tx, waypointName)
		if err != nil {
			return err
		}
		waypoint.WaypointIsSealed = false
		waypoint.WaypointSealedReason = ""
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
func UpdateWaypointTravelledCount(waypointName string) error {
	bboltWriteDb, bboltWriteDbErr := db.OpenWriteDB()
	if bboltWriteDbErr != nil {
		return bboltWriteDbErr
	}
	defer db.CloseDB(bboltWriteDb)
	return bboltWriteDb.Update(func(tx *bbolt.Tx) error {
		bucket, waypoint, err := GetWaypointForUpdate(tx, waypointName)
		if err != nil {
			return err
		}
		waypoint.WaypointTravelledCount += 1
		return PutWaypoint(bucket, waypointName, waypoint)
	})
}

// ----------------------------------
//
//	Reads every entry in the waypoint bucket and returns a slice of
//	WaypointInfo records. Uses a read-only View transaction; path existence is
//	checked for each entry inside the View and any waypoints whose paths are
//	gone are queued for sealing. The read connection is closed immediately
//	after View completes so it does not block any concurrent write connection.
//	Pending seals and corrupted-proto entries are each persisted in their own
//	follow-up Update transactions after the View closes.
//
// ----------------------------------
func GetAllWaypointsInfo() ([]WaypointInfo, error) {
	var waypointsInfo []WaypointInfo
	var corruptedWaypointName []string     // for waypoint corrupted data
	var corruptedRuneWaypointPath []string // for rune corrupted data
	waypointCorrupted := false
	runeCorrupted := false

	// seal updates are collected during the read-only View and applied afterwards;
	// calling a write transaction (Update) inside a View callback deadlocks bbolt
	type pendingSeal struct {
		name   string
		reason string
	}
	var toSeal []pendingSeal

	// open DB so we can read waypoints data
	viewErr := func() error {
		bboltReadDb, bboltReadDbErr := db.OpenReadDB()
		if bboltReadDbErr != nil {
			return bboltReadDbErr
		}
		defer db.CloseDB(bboltReadDb)

		return bboltReadDb.View(func(tx *bbolt.Tx) error {
			// ensure the waypoint bucket exists before iterating
			waypointBucket := tx.Bucket(db.WaypointBucket)
			if waypointBucket == nil {
				return fmt.Errorf("%s", style.RenderStringWithColor(i18n.LANGUAGEMAPPING.WaypointBucketNotFoundError, style.ColorError, false))
			}

			// walk every key-value pair in the bucket
			retrieveError := waypointBucket.ForEach(func(k, v []byte) error {
				// deserialize the stored proto; capture the name, set the flag, and
				// return a sentinel error to stop ForEach — recording is deferred to
				// a separate Update after the View transaction completes
				existingWaypoint := &pb.Waypoint{}
				protoErr := proto.Unmarshal(v, existingWaypoint)

				// skip corrupted data
				if protoErr != nil {
					waypointCorrupted = true
					corruptedWaypointName = append(corruptedWaypointName, string(k))
					return nil
				}

				// verify the path still exists on disk; if not, mark for sealing after View closes
				isPathExist, isPathExistErr := utils.CheckIsPathExist(existingWaypoint.WaypointPath)
				if !isPathExist {
					existingWaypoint.WaypointIsSealed = true
					existingWaypoint.WaypointSealedReason = isPathExistErr.Error()
					toSeal = append(toSeal, pendingSeal{name: string(k), reason: existingWaypoint.WaypointSealedReason})
				}

				_, runes, isCorrupted, runeErr := RetrieveRuneBasedOnWaypointPath(tx, existingWaypoint.WaypointPath)
				if isCorrupted {
					runeCorrupted = true
					corruptedRuneWaypointPath = append(corruptedRuneWaypointPath, existingWaypoint.WaypointPath)
				}

				if runeErr != nil && !isCorrupted {
					return runeErr
				}

				// construct the waypoint info type
				info := WaypointInfo{
					WaypointName:           string(k),
					WaypointPath:           existingWaypoint.WaypointPath,
					WaypointAddedAt:        existingWaypoint.WaypointAddedAt,
					WaypointTravelledCount: existingWaypoint.WaypointTravelledCount,
					WaypointIsSealed:       existingWaypoint.WaypointIsSealed,
					WaypointSealedReason:   existingWaypoint.WaypointSealedReason,
					EnterRune:              runes.EnterRunes,
					LeaveRune:              runes.LeaveRunes,
				}

				waypointsInfo = append(waypointsInfo, info)

				return nil
			})

			// wrap ForEach failure into a user-facing message
			if retrieveError != nil {
				return fmt.Errorf("%s", style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.RiftWaypointRetrieveAllError, retrieveError.Error()), style.ColorError, false))
			}

			return nil
		})
	}()

	// best-effort: persist each seal to the DB; failures are silently ignored —
	// the in-memory waypointInfo records already carry WaypointIsSealed=true
	// so the UI reflects the correct sealed state regardless
	for _, s := range toSeal {
		UpdateWaypointIsSeal(s.name, true, s.reason)
	}

	if waypointCorrupted {
		viewErr = RecordCorruptedWaypointInfo(corruptedWaypointName)
	}

	if runeCorrupted {
		viewErr = RecordCorruptedRuneInfo(corruptedRuneWaypointPath)
	}

	return waypointsInfo, viewErr
}

// ----------------------------------
//
//	Looks up a waypoint by name and returns the full deserialized *pb.Waypoint
//	record using a read-only View transaction. Returns an error when the waypoint
//	bucket is missing, the waypoint does not exist, or the stored proto is corrupted.
//
// ----------------------------------
func RetrieveWaypointInfo(waypointName string) (WaypointInfo, error) {
	waypoint := &pb.Waypoint{}
	var rune *pb.Rune
	bboltReadDb, bboltReadDbErr := db.OpenReadDB()
	var waypointInfo WaypointInfo
	if bboltReadDbErr != nil {
		return waypointInfo, bboltReadDbErr
	}
	defer db.CloseDB(bboltReadDb)

	viewErr := bboltReadDb.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(db.WaypointBucket)
		if bucket == nil {
			return fmt.Errorf("%s", style.RenderStringWithColor(i18n.LANGUAGEMAPPING.WaypointBucketNotFoundError, style.ColorError, false))
		}

		existing := bucket.Get([]byte(waypointName))
		if existing == nil {
			return fmt.Errorf("%s", style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.RiftWaypointDoNotExistsError, waypointName), style.ColorError, false))
		}

		if err := proto.Unmarshal(existing, waypoint); err != nil {
			return fmt.Errorf("%s", style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.WaypointDataCorruptedError, waypointName), style.ColorError, false))
		}

		_, retrievedRune, _, runeErr := RetrieveRuneBasedOnWaypointPath(tx, waypoint.WaypointPath)

		rune = retrievedRune

		return runeErr
	})

	if viewErr == nil {
		waypointInfo = WaypointInfo{
			WaypointName:           waypoint.WaypointName,
			WaypointPath:           waypoint.WaypointPath,
			WaypointAddedAt:        waypoint.WaypointAddedAt,
			WaypointTravelledCount: waypoint.WaypointTravelledCount,
			WaypointIsSealed:       waypoint.WaypointIsSealed,
			WaypointSealedReason:   waypoint.WaypointSealedReason,
			EnterRune:              rune.EnterRunes,
			LeaveRune:              rune.LeaveRunes,
		}
	}

	return waypointInfo, viewErr
}
