package learn

import (
	"fmt"
	"strings"

	"github.com/gohyuhan/rift/db"
	"github.com/gohyuhan/rift/i18n"
	"github.com/gohyuhan/rift/logger"
	"github.com/gohyuhan/rift/style"
	"github.com/spf13/cobra"
	"mvdan.cc/sh/v3/shell"
)

// ----------------------------------
//
//	cobra handler for the rift learn command; opens the database and delegates
//	validation and persistence to SaveSpell, then logs a success message to the
//	terminal — using the "updated" message when the spell already existed, or
//	the "learned" message when it is new
//
// ----------------------------------
var RiftLearnFunc = func(command *cobra.Command, args []string) error {
	spellName := strings.TrimSpace(args[0])
	spellCmd := strings.TrimSpace(args[1])

	// open DB and persist the new spell
	bboltDB, bboltDBErr := db.OpenDB()
	if bboltDBErr != nil {
		return bboltDBErr
	}
	defer db.CloseDB(bboltDB)

	spellExist, saveSpellErr := SaveSpell(bboltDB, spellName, spellCmd)

	if saveSpellErr != nil {
		return saveSpellErr
	}

	var message string
	var spellCmdString string
	// parse the clean cmd array
	parsedSpellCmdArray, parsedSpellCmdArrayErr := shell.Fields(spellCmd, nil)
	if parsedSpellCmdArrayErr != nil {
		spellCmdString = spellCmd
	} else {
		spellCmdString = strings.Join(parsedSpellCmdArray, " ")
	}

	if spellExist {
		message = style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.RiftSpellUpdated, spellName, spellCmdString), style.ColorGreenSoft, false)
	} else {
		message = style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.RiftSpellSaved, spellName, spellCmdString), style.ColorGreenSoft, false)
	}
	logger.LOGGER.LogToTerminal([]string{message})

	return nil
}
