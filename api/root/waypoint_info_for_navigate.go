package root

import (
	"fmt"

	apiUtils "github.com/gohyuhan/rift/api/utils"
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
func retrieveWaypointInfoForNavigate(bboltDb *bbolt.DB, waypointName string) (string, error) {
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
		viewErr = apiUtils.RecordCorruptedWaypointInfo(bboltDb, []string{waypointName})
	}

	// best-effort: seal the waypoint; failure is silently ignored — the
	// navigation error is already captured in viewErr and returned to the caller
	if needToSealWaypoint {
		apiUtils.UpdateWaypointIsSeal(bboltDb, waypointName, true, needToSealReason)
	}

	return retrievedPath, viewErr
}
