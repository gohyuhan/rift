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
	triggerWaypointRune(RUNE_ON_LEAVE)
	// Only this line goes to stdout — the shell wrapper evals it.
	fmt.Printf("cd %q", retrievedPath)

	// best-effort: increment travel count; failure is silently ignored
	UpdateWaypointTravelledCount(waypointName)
	triggerWaypointRune(RUNE_ON_ENTER)
}

// ----------------------------------
//
//	Looks up and executes the rune commands for the current working directory.
//	runeType selects RUNE_ON_ENTER or RUNE_ON_LEAVE commands. Silently returns
//	if CWD cannot be retrieved, or if no rune is registered for the path.
//	Each command is logged before execution; if CWD retrieval fails mid-loop,
//	an error is logged and the command is skipped.
//
// ----------------------------------
func triggerWaypointRune(runeType string) {
	cwd, err := utils.GetCWD()
	if err != nil {
		return
	}

	hasRune, rune := RetrieveRuneForTrigger(strings.TrimSpace(cwd))

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

	padding := strings.Repeat("  ", runeDepth)
	runeCmdsCount := len(runeCmds)
	for index, cmd := range runeCmds {
		cwd, cwdErr := utils.GetCWD()
		msg := padding + style.RenderStringWithColor(fmt.Sprintf("[%s (%v/%v) - %s]", runeType, index+1, runeCmdsCount, strings.Join(cmd.Commands, " ")), style.ColorPurpleVibrant, false)
		if cwdErr != nil {
			msg = padding + style.RenderStringWithColor(i18n.LANGUAGEMAPPING.SkippingDueToCwdErr, style.ColorError, false)
		}
		logger.LOGGER.LogToTerminal([]string{msg})
		executor.CmdExecutor().RunCmd(cmd.Commands, cwd, []string{fmt.Sprintf("RIFT_RUNE_DEPTH=%v", runeDepth+1)})

	}
}
