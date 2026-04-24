package utils

import (
	"fmt"
	"image/color"
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
//	Fires the on-leave rune for the current directory (skipped if CWD cannot
//	be resolved), increments the waypoint travel count (failure silently
//	ignored so navigation is never blocked), fires the on-enter rune for the
//	destination, then emits a cd command to stdout for the shell wrapper to eval.
//
// ----------------------------------
func ChangeDir(retrievedPath, waypointName string) {
	path, cwdErr := utils.GetCWD()
	if cwdErr == nil {
		triggerWaypointRune(RUNE_ON_LEAVE, path)
	}

	// best-effort: increment travel count; failure is silently ignored
	UpdateWaypointTravelledCount(waypointName)
	triggerWaypointRune(RUNE_ON_ENTER, retrievedPath)

	// Only this line goes to stdout — the shell wrapper evals it.
	fmt.Printf("cd %q", retrievedPath)
}

// ----------------------------------
//
//	Executes the rune commands registered for path. runeType selects
//	RUNE_ON_ENTER or RUNE_ON_LEAVE commands. Silently returns if no rune is
//	registered for path. path is passed in by the caller (not resolved
//	internally). Each command runs with path as its working directory so that
//	nested rift calls inherit the correct CWD for their own triggers.
//	Commands with an empty slice are skipped with an error log; rift
//	navigation commands are blocked with an error log.
//
// ----------------------------------
func triggerWaypointRune(runeType string, path string) {
	hasRune, rune := RetrieveRuneForTrigger(path)

	if !hasRune {
		return
	}

	runeDepth := 0
	if val, ok := os.LookupEnv("RIFT_RUNE_DEPTH"); ok {
		if n, err := strconv.Atoi(val); err == nil {
			runeDepth = n
		}
	}

	var runeCmds []*pb.RuneCmds
	var logColor color.Color
	switch runeType {
	case RUNE_ON_ENTER:
		runeCmds = rune.EnterRunes
		logColor = style.EnterRuneColorCycle[runeDepth%len(style.EnterRuneColorCycle)]
	case RUNE_ON_LEAVE:
		runeCmds = rune.LeaveRunes
		logColor = style.ColorYellowSoft
	}

	padding := strings.Repeat("  ", runeDepth)
	runeCmdsCount := len(runeCmds)
	for index, cmd := range runeCmds {
		if len(cmd.Commands) == 0 {
			errMsg := padding + style.RenderStringWithColor(fmt.Sprintf("[%s (%v/%v) - %s]", runeType, index+1, runeCmdsCount, i18n.LANGUAGEMAPPING.SkippingDueToExecutorErr), style.ColorError, false)
			logger.LOGGER.LogToTerminal([]string{errMsg})
			continue
		}
		if utils.IsRiftNavigationCommand(cmd.Commands) {
			errMsg := padding + style.RenderStringWithColor(fmt.Sprintf("[%s (%v/%v) - %s]", runeType, index+1, runeCmdsCount, i18n.LANGUAGEMAPPING.ForbiddenRiftNavigationRuneCommand), style.ColorError, false)
			logger.LOGGER.LogToTerminal([]string{errMsg})
			continue
		}
		msg := padding + style.RenderStringWithColor(fmt.Sprintf("[%s (%v/%v) - %s]", runeType, index+1, runeCmdsCount, strings.Join(cmd.Commands, " ")), logColor, false)
		logger.LOGGER.LogToTerminal([]string{msg})
		executor.CmdExecutor().ExecWithPadding(cmd.Commands, path, nil, padding)
	}
}
