package ritual

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gohyuhan/rift/executor"
	"github.com/gohyuhan/rift/i18n"
	"github.com/gohyuhan/rift/logger"
	"github.com/gohyuhan/rift/style"
	"github.com/gohyuhan/rift/utils"
	"github.com/spf13/cobra"
)

// ----------------------------------
//
//	Cobra handler for the ritual command.
//	Dispatches to ForgetRitual when --forget is passed, otherwise resolves
//	the current working directory and invokes the named ritual via InvokeRitual.
//
// ----------------------------------
var RiftRitualFunc = func(command *cobra.Command, args []string) error {
	ritualName := strings.TrimSpace(args[0])
	// --forget and casting are mutually exclusive; check which path to take
	forget, forgetErr := command.Flags().GetBool("forget")
	if forgetErr != nil {
		return fmt.Errorf("%s", style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.RiftFlagRetrieveError, "forget", forgetErr.Error()), style.ColorError, false))
	}

	if forget {
		return ForgetRitual(ritualName, true)
	} else {
		executionPath, executionPathErr := utils.GetCWD()
		if executionPathErr != nil {
			return executionPathErr
		}
		return InvokeRitual(ritualName, executionPath)
	}
}

func InvokeRitual(ritualName, executionPath string) error {
	// look up the ritual commands; errors here are user-visible (missing, corrupted)
	retrievedRitualCmds, retrieveErr := retrieveRitualInfoForInvoke(ritualName)
	if retrieveErr != nil {
		return retrieveErr
	}

	executionDepth := 0
	if val, ok := os.LookupEnv("RIFT_EXECUTION_DEPTH"); ok {
		if n, err := strconv.Atoi(val); err == nil {
			executionDepth = n
		}
	}

	logColor := style.ExecutionDepthColorCycle[executionDepth%len(style.ExecutionDepthColorCycle)]

	padding := strings.Repeat(" ", executionDepth)
	ritualCmdsCount := len(retrievedRitualCmds)
	for index, cmd := range retrievedRitualCmds {
		if len(cmd.Commands) == 0 {
			errMsg := padding + style.RenderStringWithColor(fmt.Sprintf("[RITUAL (%v/%v) - %s]", index+1, ritualCmdsCount, i18n.LANGUAGEMAPPING.RitualCommandEmpty), style.ColorError, false)
			logger.LOGGER.LogToTerminal([]string{errMsg})
			continue
		}
		if utils.IsRiftNavigationCommand(cmd.Commands) {
			errMsg := padding + style.RenderStringWithColor(fmt.Sprintf("[RITUAL (%v/%v) - %s]", index+1, ritualCmdsCount, i18n.LANGUAGEMAPPING.ForbiddenRiftNavigationRitualCommand), style.ColorError, false)
			logger.LOGGER.LogToTerminal([]string{errMsg})
			continue
		}
		msg := padding + style.RenderStringWithColor(fmt.Sprintf("[RITUAL (%v/%v) - %s]", index+1, ritualCmdsCount, strings.Join(cmd.Commands, " ")), logColor, false)
		logger.LOGGER.LogToTerminal([]string{msg})
		execErr := executor.CmdExecutor().ExecWithPadding(cmd.Commands, executionPath, nil, padding)
		if execErr != nil {
			return execErr
		}
	}

	return nil
}
