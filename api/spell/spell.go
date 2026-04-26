package spell

import (
	"fmt"
	"strings"

	apiUtils "github.com/gohyuhan/rift/api/utils"
	"github.com/gohyuhan/rift/executor"
	"github.com/gohyuhan/rift/i18n"
	"github.com/gohyuhan/rift/style"
	"github.com/gohyuhan/rift/utils"
	"github.com/spf13/cobra"
	"mvdan.cc/sh/v3/shell"
)

// ----------------------------------
//
//	Cobra handler for the spell command.
//	Dispatches to ForgetSpell when --forget is passed, otherwise resolves
//	the current working directory and casts the named spell via CastSpell.
//
// ----------------------------------
var RiftSpellFunc = func(cmd *cobra.Command, args []string) error {
	spellNameOrCmdString := strings.TrimSpace(args[0])
	// --forget and casting are mutually exclusive; check which path to take
	forgetFlag, forgetFlagErr := cmd.Flags().GetBool("forget")
	if forgetFlagErr != nil {
		return fmt.Errorf("%s", style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.RiftFlagRetrieveError, "forget", forgetFlagErr.Error()), style.ColorError, false))
	}

	if forgetFlag {
		return ForgetSpell(spellNameOrCmdString, true)
	} else {
		executionPath, executionPathErr := utils.GetCWD()
		if executionPathErr != nil {
			return executionPathErr
		}
		return CastSpell(spellNameOrCmdString, executionPath)
	}
}

// ----------------------------------
//
//	Casts a spell by name or runs a raw command string at executionPath.
//	Tries to resolve spellNameOrCmdString as a saved spell first; if found,
//	increments the cast count. If not found, falls back to parsing the input
//	as a shell command string via shell.Fields and runs it directly.
//
//	Exit code is not propagated — rift is a launcher, not a validator.
//
// ----------------------------------
func CastSpell(spellNameOrCmdString string, executionPath string) error {
	var spellCmd []string
	var castSpellErr error

	needUpdateCastCount := false

	spellCmd, castSpellErr = retrieveSpellInfoForCast(spellNameOrCmdString)
	if castSpellErr == nil {
		needUpdateCastCount = true
	} else {
		spellCmd, castSpellErr = shell.Fields(spellNameOrCmdString, nil)
		if castSpellErr != nil {
			return castSpellErr
		}
	}

	if len(spellCmd) < 1 {
		return fmt.Errorf("%s", style.RenderStringWithColor(i18n.LANGUAGEMAPPING.SpellCommandEmpty, style.ColorError, false))
	}

	spellCmdExecutor := executor.CmdExecutor().RunCmd(spellCmd, executionPath, []string{})
	if spellCmdExecutor == nil {
		errMessage := style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.InvalidSpellCommandError, strings.Join(spellCmd, " ")), style.ColorError, false)
		return fmt.Errorf("%s", errMessage)
	}

	if needUpdateCastCount {
		apiUtils.UpdateSpellCastedCount(spellNameOrCmdString)
	}

	// Run the user's command; the exit code is intentionally not propagated —
	// rift is a launcher, not a validator of the command's outcome
	spellCmdExecutor.Run()

	return nil
}

// ----------------------------------
//
//	Looks up the named spell, runs its bound terminal command at executionPath,
//	then increments the cast count (best-effort; failure is silently ignored).
//
// ----------------------------------
func RetrieveAndCastSpell(spellName string, executionPath string) error {
	// look up the spell command; errors here are user-visible (missing, corrupted)
	retrievedSpellCmd, retrieveErr := retrieveSpellInfoForCast(spellName)
	if retrieveErr != nil {
		return retrieveErr
	}

	spellCmdExecutor := executor.CmdExecutor().RunCmd(retrievedSpellCmd, executionPath, []string{})
	if spellCmdExecutor == nil {
		errMessage := style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.InvalidSpellCommandError, strings.Join(retrievedSpellCmd, " ")), style.ColorError, false)
		return fmt.Errorf("%s", errMessage)
	}

	// best-effort: increment cast count; failure is silently ignored
	apiUtils.UpdateSpellCastedCount(spellName)

	// Run the user's command; the exit code is intentionally not propagated —
	// rift is a launcher, not a validator of the command's outcome
	spellCmdExecutor.Run()

	return nil
}
