package root

import (
	"fmt"
	"strings"

	"github.com/gohyuhan/rift/api/spell"
	apiUtils "github.com/gohyuhan/rift/api/utils"
	"github.com/gohyuhan/rift/constant"
	"github.com/gohyuhan/rift/i18n"
	"github.com/gohyuhan/rift/logger"
	"github.com/gohyuhan/rift/settings"
	"github.com/gohyuhan/rift/style"
	"github.com/gohyuhan/rift/updater"
	"github.com/spf13/cobra"
)

// ----------------------------------
//
//		Cobra handler for the root rift command.
//		If --update is passed, triggers an immediate update check and returns.
//		If --version is passed, prints the current app version and returns.
//
//		When called with no args, it handles settings flags (language, autoupdate,
//		download-pre-release); if no flags were passed it falls through to help.
//
//		When called with an arg, it travels to the named waypoint by printing a
//		cd command for the shell wrapper to eval, then increments the travel count.
//
//		If --cast is passed with a spell name, it casts the spell at the waypoint path
//	    instead of just traveling there.
//
// ----------------------------------
var RiftRootFunc = func(cmd *cobra.Command, args []string) error {
	if cmd.Flags().Changed("update") {
		updater.Update()
		return nil
	}

	if cmd.Flags().Changed("version") {
		message := style.RenderStringWithColor(constant.APPVERSION, style.ColorPurpleVibrant, false)
		logger.LOGGER.LogToTerminal([]string{message})
		return nil
	}

	if len(args) < 1 {
		// settings flags can only take effect when there are no waypoint args
		languageFlagCalled := cmd.Flags().Changed("language")
		autoupdateFlagCalled := cmd.Flags().Changed("autoupdate")
		downloadPreReleaseFlagCalled := cmd.Flags().Changed("download-pre-release")

		if languageFlagCalled {
			languageSetting, languageSettingErr := cmd.Flags().GetString("language")
			if languageSettingErr != nil {
				return fmt.Errorf("%s", style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.RiftFlagRetrieveError, "language", languageSettingErr.Error()), style.ColorError, false))
			}
			settings.UpdateLanguageCode(languageSetting)
		}

		if autoupdateFlagCalled {
			autoUpdateSetting, autoUpdateSettingErr := cmd.Flags().GetBool("autoupdate")
			if autoUpdateSettingErr != nil {
				return fmt.Errorf("%s", style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.RiftFlagRetrieveError, "autoupdate", autoUpdateSettingErr.Error()), style.ColorError, false))
			}
			settings.UpdateAutoUpdate(autoUpdateSetting)
		}

		if downloadPreReleaseFlagCalled {
			downloadPreReleaseSetting, downloadPreReleaseSettingErr := cmd.Flags().GetBool("download-pre-release")
			if downloadPreReleaseSettingErr != nil {
				return fmt.Errorf("%s", style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.RiftFlagRetrieveError, "download-pre-release", downloadPreReleaseSettingErr.Error()), style.ColorError, false))
			}
			settings.UpdateDownloadPreRelease(downloadPreReleaseSetting)
		}

		// no flags were passed — show help
		if !(languageFlagCalled || autoupdateFlagCalled || downloadPreReleaseFlagCalled) {
			return cmd.Help()
		}

		return nil
	}

	waypointName := strings.TrimSpace(args[0])
	castFlagCalled := cmd.Flags().Changed("cast")

	// look up the waypoint path; errors here are user-visible (sealed, missing, corrupted)
	retrievedPath, retrieveErr := retrieveWaypointInfoForNavigate(waypointName)
	if retrieveErr != nil {
		return retrieveErr
	}

	if castFlagCalled {
		// retrieve the spell name from the --cast flag; errors here are user-visible (missing, corrupted)
		castArg, castArgErr := cmd.Flags().GetString("cast")
		if castArgErr != nil {
			return fmt.Errorf("%s", style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.RiftFlagRetrieveError, "cast", castArgErr.Error()), style.ColorError, false))
		}

		return spell.CastSpell(castArg, retrievedPath)
	} else {
		apiUtils.ChangeDir(retrievedPath, waypointName)
	}

	return nil
}
