package spellbook

import (
	"strings"

	"github.com/gohyuhan/rift/api/spell"
	"github.com/gohyuhan/rift/api/spellbook/features"
	spellbookUI "github.com/gohyuhan/rift/api/spellbook/ui"
	"github.com/gohyuhan/rift/db"
	"github.com/gohyuhan/rift/logger"
	"github.com/spf13/cobra"
)

// ----------------------------------
//
//	Cobra handler for the spellbook command.
//	With no args, launches the interactive TUI for browsing spells.
//	With a spell name arg, shows detailed info for the named spell.
//
// ----------------------------------
var RiftSpellbookFunc = func(cmd *cobra.Command, args []string) error {
	// open DB so we can read spell records
	bboltDB, bboltDBErr := db.OpenDB()
	if bboltDBErr != nil {
		return bboltDBErr
	}
	defer db.CloseDB(bboltDB)

	// no args — start spellbook Interactive UI
	if len(args) < 1 {
		spellName, spellCastPath, interactiveErr := spellbookUI.RunSpellbookInteractive(bboltDB)
		if interactiveErr != nil {
			return interactiveErr
		}

		if spellName != "" && spellCastPath != "" {
			// best-effort: increment spell cast count; failure is silently ignored
			spell.RetrieveAndCastSpell(bboltDB, spellName, spellCastPath)
		}

		return nil
	}

	// spell name arg provided — show detailed info for the named spell
	// extract and normalize the spell name from the first argument
	spellName := strings.TrimSpace(args[0])

	retrieveSpellInfoDetail, retrieveSpellInfoDetailErr := features.RetrieveSpellInfoDetail(bboltDB, spellName)

	if retrieveSpellInfoDetailErr != nil {
		return retrieveSpellInfoDetailErr
	}

	logger.LOGGER.LogToTerminal(retrieveSpellInfoDetail)

	return nil
}
