package api

import (
	"fmt"

	"github.com/gohyuhan/rift/db"
	"github.com/gohyuhan/rift/i18n"
	"github.com/gohyuhan/rift/internal/shell"
	"github.com/gohyuhan/rift/logger"
	"github.com/gohyuhan/rift/style"
	"github.com/spf13/cobra"
)

var RiftAwakenFunc = func(cmd *cobra.Command, args []string) error {
	if !shell.BinaryInPath() {
		logger.LOGGER.LogToTerminal([]string{style.RenderStringWithColor(i18n.LANGUAGEMAPPING.BinaryNotInPath, style.ColorYellowWarm, false)})
	}

	sh := shell.Detect()
	msg := style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.RiftDetectedShell, sh), style.ColorPurpleVibrant, false)
	logger.LOGGER.LogToTerminal([]string{msg})

	cfgFile, cfgFileErr := shell.ConfigFile(sh)
	if cfgFileErr != nil {
		// CMD or unsupported shell — Install returns the descriptive error
		return shell.Install(sh)
	}

	installed, isInstalledErr := shell.IsInstalled(cfgFile)
	if isInstalledErr != nil {
		return isInstalledErr
	}

	if installed {
		logger.LOGGER.LogToTerminal([]string{style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.ShellAlreadyInstalled, cfgFile), style.ColorGreenSoft, false)})
	} else {
		if err := shell.Install(sh); err != nil {
			return err
		}
	}

	dbSetupErr := db.SetupDB()
	if dbSetupErr != nil {
		return dbSetupErr
	}

	return nil
}
