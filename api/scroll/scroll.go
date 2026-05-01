package scroll

import (
	"strings"

	"github.com/gohyuhan/rift/api/ritual"
	"github.com/gohyuhan/rift/api/scroll/features"
	scrollUI "github.com/gohyuhan/rift/api/scroll/ui"
	"github.com/gohyuhan/rift/logger"
	"github.com/spf13/cobra"
)

// ----------------------------------
//
//	Cobra handler for the scroll command.
//	With no args, launches the interactive TUI for browsing rituals.
//	With a ritual name arg, shows detailed info for the named ritual.
//
// ----------------------------------
var RiftScrollFunc = func(cmd *cobra.Command, args []string) error {
	// no args — start spellbook Interactive UI
	if len(args) < 1 {
		ritualName, ritualInvokePath, interactiveErr := scrollUI.RunScrollInteractive()
		if interactiveErr != nil {
			return interactiveErr
		}

		if ritualName != "" && ritualInvokePath != "" {
			// Cast the spell; the cast count update is best-effort, but execution errors are returned.
			return ritual.InvokeRitual(ritualName, ritualInvokePath)
		}

		return nil
	}

	// ritual name arg provided — show detailed info for the named ritual
	// extract and normalize the ritual name from the first argument
	ritualName := strings.TrimSpace(args[0])

	retrieveRitualInfoDetail, retrieveRitualInfoDetailErr := features.RetrieveRitualInfoDetail(ritualName)

	if retrieveRitualInfoDetailErr != nil {
		return retrieveRitualInfoDetailErr
	}

	logger.LOGGER.LogToTerminal(retrieveRitualInfoDetail)

	return nil
}
