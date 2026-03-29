package spell

import (
	"fmt"
	"strings"

	apiUtils "github.com/gohyuhan/rift/api/utils"
	"github.com/gohyuhan/rift/db"
	"github.com/gohyuhan/rift/executor"
	"github.com/gohyuhan/rift/i18n"
	"github.com/gohyuhan/rift/style"
	"github.com/gohyuhan/rift/utils"
	"github.com/spf13/cobra"
	"go.etcd.io/bbolt"
)

// ----------------------------------
//
//	Cobra handler for the spell command.
//	Looks up the named spell in the DB, executes its bound terminal command,
//	then increments the spell's cast count.
//
// ----------------------------------
var RiftSpellFunc = func(cmd *cobra.Command, args []string) error {
	spellName := strings.TrimSpace(args[0])

	// open DB for reading spell data
	bboltDB, bboltDBErr := db.OpenDB()
	if bboltDBErr != nil {
		return bboltDBErr
	}
	defer db.CloseDB(bboltDB)

	executionPath, executionPathErr := utils.GetCWD()
	if executionPathErr != nil {
		return executionPathErr
	}
	return RetrieveAndCastSpell(bboltDB, spellName, executionPath)
}

// ----------------------------------
//
//	Looks up the named spell, runs its bound terminal command at executionPath,
//	then increments the cast count (best-effort; failure is silently ignored).
//
// ----------------------------------
func RetrieveAndCastSpell(bboltDB *bbolt.DB, spellName string, executionPath string) error {
	// look up the spell command; errors here are user-visible (missing, corrupted)
	retrievedSpellCmd, retrieveErr := retrieveSpellInfoForCast(bboltDB, spellName)
	if retrieveErr != nil {
		return retrieveErr
	}

	spellCmdExecutor := executor.CmdExecutor().RunCmd(retrievedSpellCmd, executionPath)
	if spellCmdExecutor == nil {
		errMessage := style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.InvalidSpellCommandError, strings.Join(retrievedSpellCmd, " ")), style.ColorError, false)
		return fmt.Errorf("%s", errMessage)
	}
	// Run the user's command; the exit code is intentionally not propagated —
	// rift is a launcher, not a validator of the command's outcome
	spellCmdExecutor.Run()

	// best-effort: increment cast count; failure is silently ignored
	apiUtils.UpdateSpellCastedCount(bboltDB, spellName)

	return nil
}
