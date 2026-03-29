// ----------------------------------
//
//	SPELL RELATED UTILS
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
//	Persists the spell names into the corrupted-records bucket via a fresh
//	Update transaction (independent of any prior transaction that may have
//	failed), then returns a user-facing corruption error.
//	The Update's own error is intentionally not propagated — recording is
//	best-effort; the corruption error is always returned to the caller.
//
// ----------------------------------
func RecordCorruptedSpellInfo(bboltDB *bbolt.DB, corruptedSpellsName []string) error {
	// best-effort write — ignore the Update error; the caller always gets the corruption message
	bboltDB.Update(func(tx *bbolt.Tx) error {
		spellCorruptedBucket := tx.Bucket(db.SpellDataCorruptedBucketRecord)
		if spellCorruptedBucket != nil {
			for _, corruptedSpell := range corruptedSpellsName {
				spellCorruptedBucket.Put([]byte(corruptedSpell), []byte(corruptedSpell))
			}
		}
		return nil
	})
	return fmt.Errorf("%s", style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.SpellDataCorruptedError, strings.Join(corruptedSpellsName, ",")), style.ColorError, false))
}

// ----------------------------------
//
//	Fetches and deserializes the named spell from the bucket within an
//	already-open Update transaction. Returns the bucket, the deserialized
//	record, or an error if the bucket is missing, the spell does not exist,
//	or the stored proto is corrupted. Callers mutate the returned record and
//	re-persist it via bucket.Put.
//
// ----------------------------------
func GetSpellForUpdate(tx *bbolt.Tx, spellName string) (*bbolt.Bucket, *pb.Spell, error) {
	bucket := tx.Bucket(db.SpellBucket)
	if bucket == nil {
		return nil, nil, fmt.Errorf("%s", style.RenderStringWithColor(i18n.LANGUAGEMAPPING.SpellBucketNotFoundError, style.ColorError, false))
	}

	existing := bucket.Get([]byte(spellName))
	if existing == nil {
		return nil, nil, fmt.Errorf("%s", style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.RiftSpellDoNotExistsError, spellName), style.ColorError, false))
	}

	spell := &pb.Spell{}
	if err := proto.Unmarshal(existing, spell); err != nil {
		return nil, nil, fmt.Errorf("%s", style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.SpellDataCorruptedError, spellName), style.ColorError, false))
	}

	return bucket, spell, nil
}

// ----------------------------------
//
//	Persists a mutated spell record back into its bucket. Returns an error
//	if marshalling or the bucket write fails.
//
// ----------------------------------
func PutSpell(bucket *bbolt.Bucket, spellName string, spell *pb.Spell) error {
	data, err := proto.Marshal(spell)
	if err != nil {
		return fmt.Errorf("%s", style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.RiftSpellUpdateError, spellName, err.Error()), style.ColorError, false))
	}
	return bucket.Put([]byte(spellName), data)
}

// ----------------------------------
//
//	Increments the cast count for the named spell in the DB bucket.
//	Returns an error if the bucket is missing, the data is unreadable, or the
//	spell does not exist.
//
// ----------------------------------
func UpdateSpellCastedCount(bboltDb *bbolt.DB, spellName string) error {
	return bboltDb.Update(func(tx *bbolt.Tx) error {
		bucket, spell, err := GetSpellForUpdate(tx, spellName)
		if err != nil {
			return err
		}
		spell.SpellCastCount += 1
		return PutSpell(bucket, spellName, spell)
	})
}
