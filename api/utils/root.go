package utils

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gohyuhan/rift/executor"
	"github.com/gohyuhan/rift/i18n"
	"github.com/gohyuhan/rift/logger"
	pb "github.com/gohyuhan/rift/proto"
	"github.com/gohyuhan/rift/style"
	"github.com/gohyuhan/rift/utils"
)

const (
	RUNE_ON_ENTER = "RUNE_ON_ENTER"
	RUNE_ON_LEAVE = "RUNE_ON_LEAVE"
)

// ----------------------------------
//
//	Fires the on-leave rune for the current directory, emits a cd command to
//	stdout for the shell wrapper to eval, increments the waypoint travel count,
//	then fires the on-enter rune for the destination. Travel-count failure is
//	silently ignored so navigation is never blocked.
//
// ----------------------------------
func ChangeDir(retrievedPath, waypointName string) {
	leavePath, cwdErr := utils.GetCWD()
	if cwdErr == nil {
		triggerWaypointRune(RUNE_ON_LEAVE, strings.TrimSpace(leavePath))
	}

	// best-effort: increment travel count; failure is silently ignored
	UpdateWaypointTravelledCount(waypointName)
	triggerWaypointRune(RUNE_ON_ENTER, retrievedPath)

	// Only this line goes to stdout — the shell wrapper evals it.
	// Skip at depth>0: nested rift calls run inside executor subprocesses;
	// their stdout is not eval'd by the shell wrapper.
	if val, ok := os.LookupEnv("RIFT_RUNE_DEPTH"); !ok || val == "0" {
		fmt.Printf("cd %q", retrievedPath)
	}
}

// ----------------------------------
//
//	Looks up and executes the rune commands for path. runeType selects
//	RUNE_ON_ENTER or RUNE_ON_LEAVE commands. Silently returns if no rune is
//	registered for path. Each command runs with path as its working directory
//	so that nested rift calls inherit the correct CWD for their own triggers.
//
// ----------------------------------
func triggerWaypointRune(runeType string, path string) {
	hasRune, rune := RetrieveRuneForTrigger(path)

	if !hasRune {
		return
	}

	var runeCmds []*pb.RuneCmds
	switch runeType {
	case RUNE_ON_ENTER:
		runeCmds = rune.EnterRunes
	case RUNE_ON_LEAVE:
		runeCmds = rune.LeaveRunes
	}

	runeDepth := 0
	if val, ok := os.LookupEnv("RIFT_RUNE_DEPTH"); ok {
		if n, err := strconv.Atoi(val); err == nil {
			runeDepth = n
		}
	}

	// at depth>0, no real chdir has occurred — outer call already fired LEAVE
	if runeType == RUNE_ON_LEAVE && runeDepth > 0 {
		return
	}

	padding := strings.Repeat("  ", runeDepth)
	runeCmdsCount := len(runeCmds)
	for index, cmd := range runeCmds {
		runeExecutor := executor.CmdExecutor().RunCmd(cmd.Commands, path, []string{fmt.Sprintf("RIFT_RUNE_DEPTH=%v", runeDepth+1)})
		if runeExecutor != nil {
			msg := padding + style.RenderStringWithColor(fmt.Sprintf("[%s (%v/%v) - %s]", runeType, index+1, runeCmdsCount, strings.Join(cmd.Commands, " ")), style.ColorPurpleVibrant, false)
			logger.LOGGER.LogToTerminal([]string{msg})
			runeExecutor.Run()
		} else {
			errMsg := padding + style.RenderStringWithColor(fmt.Sprintf("[%s (%v/%v) - %s]", runeType, index+1, runeCmdsCount, i18n.LANGUAGEMAPPING.SkippingDueToExecutorErr), style.ColorError, false)
			logger.LOGGER.LogToTerminal([]string{errMsg})
		}
	}
}
