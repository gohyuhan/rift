package learn

import (
	"fmt"
	"slices"
	"time"

	apiUtils "github.com/gohyuhan/rift/api/utils"
	"github.com/gohyuhan/rift/constant"
	"github.com/gohyuhan/rift/db"
	"github.com/gohyuhan/rift/i18n"
	pb "github.com/gohyuhan/rift/proto"
	"github.com/gohyuhan/rift/style"
	"github.com/gohyuhan/rift/utils"
	"go.etcd.io/bbolt"
	"google.golang.org/protobuf/proto"
	"mvdan.cc/sh/v3/shell"
)

// ----------------------------------
//
//	Parses spellCmd into a shell-quoted field array, validates spellName against
//	reserved keywords and naming rules, and rejects empty commands and any
//	command containing cd. Persists the spell to the spell bucket in a single
//	Update transaction, overwriting any existing entry with a fresh record
//	(cast count reset to 0, added-at timestamp set to now UTC).
//	Returns true when an existing spell was overwritten.
//
// ----------------------------------
func SaveSpell(spellName string, spellCmd string) (bool, error) {
	hasExisting := false

	spellCommandArray, spellCommandArrayErr := shell.Fields(spellCmd, nil)
	if spellCommandArrayErr != nil {
		return hasExisting, spellCommandArrayErr
	}

	// reject names that clash with rift's own subcommands
	if err := apiUtils.CheckIfKeywordIsReservedForRift(spellName); err != nil {
		return hasExisting, err
	}

	// reject names that contain spaces
	if err := apiUtils.IsNickNameValid(spellName); err != nil {
		return hasExisting, err
	}

	bboltWriteDb, bboltWriteDbErr := db.OpenWriteDB()
	if bboltWriteDbErr != nil {
		return hasExisting, bboltWriteDbErr
	}
	defer db.CloseDB(bboltWriteDb)

	bboltWriteDbErr = bboltWriteDb.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(db.SpellBucket)
		if bucket == nil {
			return fmt.Errorf("%s", style.RenderStringWithColor(i18n.LANGUAGEMAPPING.SpellBucketNotFoundError, style.ColorError, false))
		}

		// check if the command is not empty
		if len(spellCommandArray) < 1 {
			errMessage := style.RenderStringWithColor(i18n.LANGUAGEMAPPING.SpellCommandEmpty, style.ColorError, false)
			return fmt.Errorf("%s", errMessage)
		}

		// check if the command starts with a shell built-in
		if slices.Contains(constant.ShellBuildInCmd, spellCommandArray[0]) {
			errMessage := style.RenderStringWithColor(i18n.LANGUAGEMAPPING.ForbiddenShellBuiltinSpellCommand, style.ColorError, false)
			return fmt.Errorf("%s", errMessage)
		}

		// check if the command is a rift waypoint navigation command
		if utils.IsRiftNavigationCommand(spellCommandArray) {
			errMessage := style.RenderStringWithColor(i18n.LANGUAGEMAPPING.ForbiddenRiftNavigationSpellCommand, style.ColorError, false)
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

	return hasExisting, bboltWriteDbErr
}
