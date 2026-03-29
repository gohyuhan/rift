package spell

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
//	Deletes the named spell from the spell bucket within a single write
//	transaction. Also removes any corrupted-data record for the spell if one
//	exists. The delete is intentionally idempotent — forgetting a non-existent
//	spell succeeds silently. Logs a success message to the terminal when
//	logToTerminal is true.
//
// ----------------------------------
func ForgetSpell(bboltDb *bbolt.DB, spellName string, logToTerminal bool) error {
	return bboltDb.Update(func(tx *bbolt.Tx) error {
		// ensure the spell bucket exists before attempting the delete
		spellBucket := tx.Bucket(db.SpellBucket)
		if spellBucket == nil {
			return fmt.Errorf("%s", style.RenderStringWithColor(i18n.LANGUAGEMAPPING.SpellBucketNotFoundError, style.ColorError, false))
		}

		// bbolt.Delete returns nil for missing keys; the behavior here is
		// intentionally idempotent — destroying a non-existent spell succeeds
		forgetSpellErr := spellBucket.Delete([]byte(spellName))
		if forgetSpellErr != nil {
			return fmt.Errorf("%s", style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.RiftSpellForgetError, spellName, forgetSpellErr.Error()), style.ColorError, false))
		}

		// when a spell was forgotten also destroy the corrupted data record for the spell (if available)
		corruptedSpellBucket := tx.Bucket(db.SpellDataCorruptedBucketRecord)
		if corruptedSpellBucket != nil {
			corruptedSpellBucket.Delete([]byte(spellName))
		}

		// report the forgetting to the terminal
		if logToTerminal {
			message := style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.RiftSpellForgetSuccess, spellName), style.ColorGreenSoft, false)
			logger.LOGGER.LogToTerminal([]string{message})
		}

		return nil
	})
}
