//	  ▄████████  ▄█     ▄████████     ███
//	 ███    ███ ███    ███    ███ ▀█████████▄
//	 ███    ███ ███▌   ███    █▀     ▀███▀▀██
//	▄███▄▄▄▄██▀ ███▌  ▄███▄▄▄         ███   ▀
// ▀▀███▀▀▀▀▀   ███▌ ▀▀███▀▀▀         ███
// ▀███████████ ███    ███            ███
//	███    ███ ███    ███            ███
//	███    ███ █▀     ███           ▄████▀
//	███    ███

package main

import (
	"fmt"

	"github.com/gohyuhan/rift/api"
	"github.com/gohyuhan/rift/cmd"
	"github.com/gohyuhan/rift/i18n"
	"github.com/gohyuhan/rift/logger"
	"github.com/gohyuhan/rift/settings"
	"github.com/gohyuhan/rift/style"
	"github.com/gohyuhan/rift/updater"
)

func main() {
	logger.InitLogger()
	settings.InitOrReadSettings()

	// check for update if user allows it
	if settings.RIFTSETTINGS.AutoUpdate {
		updater.AutoUpdater()
	}
	cARErr := api.CheckAndRunSetup()
	if cARErr != nil {
		errMsg := style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.CheckAndRunSetupError, cARErr.Error()), style.ColorError, false)
		logger.LOGGER.LogToTerminal([]string{errMsg})
	}
	cmd.Execute()
}
