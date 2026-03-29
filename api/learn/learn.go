package learn

import (
	"fmt"
	"strings"

	apiUtils "github.com/gohyuhan/rift/api/utils"
	"github.com/gohyuhan/rift/db"
	"github.com/gohyuhan/rift/i18n"
	"github.com/gohyuhan/rift/logger"
	"github.com/gohyuhan/rift/style"
	"github.com/spf13/cobra"
)

// ----------------------------------
//
//	Cobra handler for the learn command.
//	Validates the spell name against reserved keywords and naming rules,
//	then persists the spell (name → command) into the DB,
//	creating a new entry or overriding an existing one.
//
// ----------------------------------
var RiftLearnFunc = func(command *cobra.Command, args []string) error {
	spellName := strings.TrimSpace(args[0])
	spellCmd := strings.TrimSpace(args[1])
	spellCmdArray := strings.Fields(spellCmd)

	// reject names that clash with rift's own subcommands
	if err := apiUtils.CheckIfKeywordIsReservedForRift(spellName); err != nil {
		return err
	}

	// reject names that contain spaces
	if err := apiUtils.IsNickNameValid(spellName); err != nil {
		return err
	}

	// open DB and persist the new waypoint
	bboltDB, bboltDBErr := db.OpenDB()
	if bboltDBErr != nil {
		return bboltDBErr
	}
	defer db.CloseDB(bboltDB)

	spellExist, saveSpellErr := saveSpell(bboltDB, spellName, spellCmdArray)

	if saveSpellErr != nil {
		return saveSpellErr
	}
	var message string
	if spellExist {
		message = style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.RiftSpellUpdated, spellName, strings.Join(spellCmdArray, " ")), style.ColorGreenSoft, false)
	} else {
		message = style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.RiftSpellSaved, spellName, strings.Join(spellCmdArray, " ")), style.ColorGreenSoft, false)
	}
	logger.LOGGER.LogToTerminal([]string{message})

	return nil
}
