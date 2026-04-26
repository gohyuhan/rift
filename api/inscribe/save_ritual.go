package inscribe

import (
	"fmt"
	"time"

	apiUtils "github.com/gohyuhan/rift/api/utils"
	"github.com/gohyuhan/rift/db"
	"github.com/gohyuhan/rift/i18n"
	pb "github.com/gohyuhan/rift/proto"
	"github.com/gohyuhan/rift/style"
	"go.etcd.io/bbolt"
)

// ----------------------------------
//
//	Validates ritualName against reserved keywords and naming rules, normalizes
//	ritualCmdString into a command slice, and rejects empty or invalid commands.
//	Persists the ritual to the ritual bucket in a single Update transaction.
//	If a ritual with the same name exists, override=true updates its description
//	and commands; override=false returns an error.
//
// ----------------------------------
func SaveRitual(ritualName string, ritualDesc string, ritualCmdString string, override bool) error {
	if err := apiUtils.IsNickNameValid(ritualName); err != nil {
		return err
	}

	ritualCmds, ritualCmdsErr := apiUtils.NormalizeAndCheckRitualCommandsAreValid(ritualCmdString)
	if ritualCmdsErr != nil {
		return ritualCmdsErr
	}

	// fail fast before opening DB
	if len(ritualCmds) < 1 {
		return fmt.Errorf("%s", style.RenderStringWithColor(i18n.LANGUAGEMAPPING.RitualCommandEmpty, style.ColorError, false))
	}

	bboltWriteDb, bboltWriteDbErr := db.OpenWriteDB()
	if bboltWriteDbErr != nil {
		return bboltWriteDbErr
	}
	defer db.CloseDB(bboltWriteDb)

	bboltWriteDbErr = bboltWriteDb.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(db.RitualBucket)
		if bucket == nil {
			return fmt.Errorf("%s", style.RenderStringWithColor(i18n.LANGUAGEMAPPING.RitualBucketNotFoundError, style.ColorError, false))
		}

		_, ritual, _ := apiUtils.GetRitualForUpdate(tx, ritualName)

		if ritual != nil && override {
			// existing ritual — override: update desc and cmds in place
			ritual.RitualDesc = ritualDesc
			ritual.RitualCmds = ritualCmds
		} else if ritual != nil && !override {
			// existing ritual — no override flag: reject with hint to use --override
			return fmt.Errorf("%s", style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.RitualExistAndOverrideNotAllowedError, ritualName), style.ColorError, false))
		} else {
			// no existing ritual: create new
			ritual = &pb.Ritual{
				RitualName:        ritualName,
				RitualDesc:        ritualDesc,
				RitualCmds:        ritualCmds,
				RitualAddedAt:     time.Now().UTC().Format(time.RFC3339),
				RitualInvokeCount: 0,
			}
		}

		return apiUtils.PutRitual(bucket, ritualName, ritual)
	})

	return bboltWriteDbErr
}
