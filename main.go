// ██████╗ ██╗███████╗████████╗
// ██╔══██╗██║██╔════╝╚══██╔══╝
// ██████╔╝██║█████╗     ██║
// ██╔══██╗██║██╔══╝     ██║
// ██║  ██║██║██║        ██║
// ╚═╝  ╚═╝╚═╝╚═╝        ╚═╝

package main

import (
	"github.com/gohyuhan/rift/cmd"
	"github.com/gohyuhan/rift/i18n"
	"github.com/gohyuhan/rift/logger"
	"github.com/gohyuhan/rift/settings"
	"github.com/gohyuhan/rift/updater"
)

func main() {
	logger.InitLogger()
	settings.InitOrReadSettings()
	i18n.InitRiftLanguageMapping(settings.RIFTSETTINGS.LanguageCode)

	// check for update if user allows it
	if settings.RIFTSETTINGS.AutoUpdate {
		updater.AutoUpdater()
	}
	cmd.Execute()
}
