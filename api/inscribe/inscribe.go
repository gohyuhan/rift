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

	override, overrideErr := command.Flags().GetBool("override")
	if overrideErr != nil {
		return fmt.Errorf("%s", style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.RiftFlagRetrieveError, "override", overrideErr.Error()), style.ColorError, false))
	}

	saveRitualErr := SaveRitual(ritualName, ritualDesc, ritualCmds, override)

	if saveRitualErr != nil {
		return saveRitualErr
	}

	message := style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.RiftRitualSaved, ritualName), style.ColorGreenSoft, false)
	logger.LOGGER.LogToTerminal([]string{message})

	return nil
}
