package spell

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
//	Looks up a spell by name and returns its stored command.
//	Uses a read-only View transaction for the lookup; corruption recording
//	is deferred to a follow-up Update transaction — bbolt write locks are not
//	reentrant, so calling Update inside View would deadlock.
//	Returns the stored command slice, or an error when:
//	  - the spell bucket is missing
//	  - the spell does not exist
//	  - the stored proto data is corrupted
//
// ----------------------------------
func retrieveSpellInfoForCast(spellName string) ([]string, error) {
	retrievedCmd := []string{}
	spellCorrupted := false

	viewErr := func() error {
		// open DB for reading spell data
		bboltReadDb, bboltReadDbErr := db.OpenReadDB()
		if bboltReadDbErr != nil {
			return bboltReadDbErr
		}
		defer db.CloseDB(bboltReadDb)

		return bboltReadDb.View(func(tx *bbolt.Tx) error {
			bucket := tx.Bucket(db.SpellBucket)
			if bucket == nil {
				return fmt.Errorf("%s", style.RenderStringWithColor(i18n.LANGUAGEMAPPING.SpellBucketNotFoundError, style.ColorError, false))
			}

			// check the spell exists in the bucket
			existing := bucket.Get([]byte(spellName))
			if existing == nil {
				errorMessage := style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.RiftSpellDoNotExistsError, spellName), style.ColorError, false)
				return fmt.Errorf("%s", errorMessage)
			}

			// deserialize the stored proto; set the flag and return nil so the View
			// commits cleanly — corruption recording is deferred to a follow-up Update
			spellInfo := &pb.Spell{}
			protoErr := proto.Unmarshal(existing, spellInfo)
			if protoErr != nil {
				spellCorrupted = true
				return nil
			}

			retrievedCmd = spellInfo.SpellCommand
			return nil
		})
	}()

	if spellCorrupted {
		viewErr = apiUtils.RecordCorruptedSpellInfo([]string{spellName})
	}

	return retrievedCmd, viewErr
}
