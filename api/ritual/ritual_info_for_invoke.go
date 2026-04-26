package ritual

import (
	"fmt"

	apiUtils "github.com/gohyuhan/rift/api/utils"
	"github.com/gohyuhan/rift/db"
	"github.com/gohyuhan/rift/i18n"
	pb "github.com/gohyuhan/rift/proto"
	"github.com/gohyuhan/rift/style"
	"go.etcd.io/bbolt"
	"google.golang.org/protobuf/proto"
)

// ----------------------------------
//
//	Looks up a ritual by name and returns its stored commands.
//	Uses a read-only View transaction for the lookup; corruption recording
//	is deferred to a follow-up Update transaction — bbolt write locks are not
//	reentrant, so calling Update inside View would deadlock.
//	Returns the stored command slice, or an error when:
//	  - the ritual bucket is missing
//	  - the ritual does not exist
//	  - the stored proto data is corrupted
//
// ----------------------------------
func retrieveRitualInfoForInvoke(ritualName string) ([]*pb.RitualCmds, error) {
	var retrievedRitualCmds []*pb.RitualCmds
	ritualCorrupted := false

	// open DB for reading ritual data
	bboltReadDb, bboltReadDbErr := db.OpenReadDB()
	if bboltReadDbErr != nil {
		return retrievedRitualCmds, bboltReadDbErr
	}
	defer db.CloseDB(bboltReadDb)

	viewErr := bboltReadDb.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(db.RitualBucket)
		if bucket == nil {
			return fmt.Errorf("%s", style.RenderStringWithColor(i18n.LANGUAGEMAPPING.RitualBucketNotFoundError, style.ColorError, false))
		}

		// check the ritual exists in the bucket
		existing := bucket.Get([]byte(ritualName))
		if existing == nil {
			errorMessage := style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.RiftRitualDoNotExistsError, ritualName), style.ColorError, false)
			return fmt.Errorf("%s", errorMessage)
		}

		// deserialize the stored proto; set the flag and return nil so the View
		// commits cleanly — corruption recording is deferred to a follow-up Update
		ritualInfo := &pb.Ritual{}
		protoErr := proto.Unmarshal(existing, ritualInfo)
		if protoErr != nil {
			ritualCorrupted = true
			return nil
		}

		retrievedRitualCmds = ritualInfo.RitualCmds
		return nil
	})

	// close early so it will not block write connection below
	db.CloseDB(bboltReadDb)

	if ritualCorrupted {
		viewErr = apiUtils.RecordCorruptedRitualInfo([]string{ritualName})
	}

	return retrievedRitualCmds, viewErr
}
