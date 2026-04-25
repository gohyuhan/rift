package rune

import (
	"fmt"
	"strings"

	runeUI "github.com/gohyuhan/rift/api/rune/ui"
	"github.com/gohyuhan/rift/i18n"
	"github.com/gohyuhan/rift/logger"
	"github.com/spf13/cobra"
)

// ----------------------------------
//
//	Cobra handler for the rune command.
//	Launches the rune engraving interactive TUI for the given waypoint;
//	logs a success or cancellation message once the session ends.
//
// ----------------------------------
var RiftRuneFunc = func(cmd *cobra.Command, args []string) error {
	// extract and normalise the waypoint name from the first argument
	waypointName := strings.TrimSpace(args[0])

	success, interactiveErr := runeUI.RunRuneInteractive(waypointName)
	if interactiveErr != nil {
		return interactiveErr
	}

	if success {
		logger.LOGGER.LogToTerminal([]string{fmt.Sprintf(i18n.LANGUAGEMAPPING.RiftRuneEngraveSuccessful, waypointName)})
	} else {
		logger.LOGGER.LogToTerminal([]string{fmt.Sprintf(i18n.LANGUAGEMAPPING.RiftRuneEngraveNone, waypointName)})
	}

	return nil
}
