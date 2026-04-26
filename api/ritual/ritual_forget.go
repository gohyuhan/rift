package ritual

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
//	Deletes the named ritual from the ritual bucket within a single write
//	transaction. Also removes any corrupted-data record for the ritual if one
//	exists. The delete is intentionally idempotent — forgetting a non-existent
//	ritual succeeds silently. Logs a success message to the terminal when
//	logToTerminal is true.
//
// ----------------------------------
func ForgetRitual(ritualName string, logToTerminal bool) error {
	bboltWriteDb, bboltWriteDbErr := db.OpenWriteDB()
	if bboltWriteDbErr != nil {
		return bboltWriteDbErr
	}
	defer db.CloseDB(bboltWriteDb)

	dbErr := bboltWriteDb.Update(func(tx *bbolt.Tx) error {
		// ensure the ritual bucket exists before attempting the delete
		ritualBucket := tx.Bucket(db.RitualBucket)
		if ritualBucket == nil {
			return fmt.Errorf("%s", style.RenderStringWithColor(i18n.LANGUAGEMAPPING.RitualBucketNotFoundError, style.ColorError, false))
		}

		// bbolt.Delete returns nil for missing keys — forgetting a non-existent ritual succeeds silently
		forgetRitualErr := ritualBucket.Delete([]byte(ritualName))
		if forgetRitualErr != nil {
			return fmt.Errorf("%s", style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.RiftRitualForgetError, ritualName, forgetRitualErr.Error()), style.ColorError, false))
		}

		// also remove the corrupted-data record for this ritual if one exists
		corruptedRitualBucket := tx.Bucket(db.RitualDataCorruptedBucketRecord)
		if corruptedRitualBucket != nil {
			corruptedRitualBucket.Delete([]byte(ritualName))
		}

		return nil
	})

	// log success to terminal
	if dbErr == nil && logToTerminal {
		message := style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.RiftRitualForgetSuccess, ritualName), style.ColorGreenSoft, false)
		logger.LOGGER.LogToTerminal([]string{message})
	}
	return dbErr
}
