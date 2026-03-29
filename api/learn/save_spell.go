package learn

import (
	"fmt"
	"slices"
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
//	Persists a new spell (name → path) into the spell bucket.
//	learn on existing spell name will overwrite the command, reset the cast count, and update the added timestamp to now.
//
// ----------------------------------
func saveSpell(bboltDb *bbolt.DB, spellName string, spellCommandArray []string) (bool, error) {
	hasExisting := false
	dbErr := bboltDb.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(db.SpellBucket)
		if bucket == nil {
			return fmt.Errorf("%s", style.RenderStringWithColor(i18n.LANGUAGEMAPPING.SpellBucketNotFoundError, style.ColorError, false))
		}

		if slices.Contains(spellCommandArray, "cd") {
			errMessage := style.RenderStringWithColor(i18n.LANGUAGEMAPPING.ForbiddenCDSpellCommand, style.ColorError, false)
			return fmt.Errorf("%s", errMessage)
		}

		existing := bucket.Get([]byte(spellName))
		if existing != nil {
			hasExisting = true
		}

		// build the new spell record with defaults
		spell := &pb.Spell{
			SpellName:      spellName,
			SpellCommand:   spellCommandArray,
			SpellAddedAt:   time.Now().UTC().Format(time.RFC3339),
			SpellCastCount: 0,
		}

		data, err := proto.Marshal(spell)
		if err != nil {
			return err
		}

		return bucket.Put([]byte(spellName), data)
	})

	return hasExisting, dbErr
}
