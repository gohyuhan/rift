package inscribe

import (
	"fmt"
	"strings"

	"github.com/gohyuhan/rift/i18n"
	"github.com/gohyuhan/rift/logger"
	"github.com/gohyuhan/rift/style"
	"github.com/spf13/cobra"
)

var RiftInscribeFunc = func(command *cobra.Command, args []string) error {
	ritualName := strings.TrimSpace(args[0])
	ritualDesc := strings.TrimSpace(args[1])
	ritualCmds := strings.TrimSpace(args[2])

	overrideFlagCalled := command.Flags().Changed("override")

	saveRitualErr := SaveRitual(ritualName, ritualDesc, ritualCmds, overrideFlagCalled)

	if saveRitualErr != nil {
		return saveRitualErr
	}

	message := style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.RiftRitualSaved, ritualName), style.ColorGreenSoft, false)
	logger.LOGGER.LogToTerminal([]string{message})

	return nil
}
