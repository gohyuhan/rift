// -----------------
//
//	GENERAL API RELATED UTILS
//
// ----------------------------------

package utils

import (
	"fmt"
	"os"
	"slices"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/gohyuhan/rift/constant"
	"github.com/gohyuhan/rift/db"
	"github.com/gohyuhan/rift/i18n"
	"github.com/gohyuhan/rift/internal/shell"
	"github.com/gohyuhan/rift/logger"
	"github.com/gohyuhan/rift/settings"
	"github.com/gohyuhan/rift/style"
	"github.com/gohyuhan/rift/utils"
	"github.com/spf13/cobra"
	"golang.org/x/mod/semver"
)

// ----------------------------------
//
//	Runs the full rift setup sequence:
//	  1. Warns if the rift binary is not on PATH
//	  2. Detects the current shell
//	  3. Installs the shell integration into the shell config file (if not already present)
//	  4. Initializes (or migrates) the bbolt DB
//
// ----------------------------------
func RiftSetup() error {
	// warn if the binary won't be found after this session ends
	if !shell.BinaryInPath() {
		message := style.RenderStringWithColor(i18n.LANGUAGEMAPPING.BinaryNotInPath, style.ColorYellowWarm, false)
		logger.LOGGER.LogToTerminal([]string{message})
	}

	// detect the running shell so we know which config file to modify
	sh := shell.Detect()
	msg := style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.RiftDetectedShell, sh), style.ColorPurpleVibrant, false)
	logger.LOGGER.LogToTerminal([]string{msg})

	cfgFile, cfgFileErr := shell.ConfigFile(sh)
	if cfgFileErr != nil {
		// CMD or unsupported shell — Install returns the descriptive error
		return shell.Install(sh)
	}

	// check whether the shell integration snippet is already present
	installed, isInstalledErr := shell.IsInstalled(cfgFile)
	if isInstalledErr != nil {
		return isInstalledErr
	}

	if installed {
		var messages []string
		messages = append(messages, style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.ShellAlreadyInstalled, cfgFile), style.ColorGreenSoft, false))
		messages = append(messages, style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.ShellReloadHint, shell.ReloadHint(sh, cfgFile)), style.ColorCyanSoft, false))
		logger.LOGGER.LogToTerminal(messages)
	} else {
		if err := shell.Install(sh); err != nil {
			return err
		}
	}

	// just to change to new line
	logger.LOGGER.LogToTerminal([]string{""})

	// initialize the DB (creates buckets if they don't exist yet)
	dbSetupErr := db.SetupDB()
	if dbSetupErr != nil {
		return dbSetupErr
	}

	return nil
}

// ----------------------------------
//
//	Decides whether setup needs to run and triggers it if so.
//	Setup is required when any of the following is true:
//	  - the settings file does not exist on disk
//	  - the DB file does not exist on disk
//	  - the binary version is newer than the version recorded in settings
//	After a successful setup the settings version is stamped with the current binary version.
//
// ----------------------------------
func CheckAndRunSetup() error {
	needSetup := false

	// check settings file presence
	settingsPath, err := utils.GetRiftSettingsFilePath()
	if err != nil {
		return err
	}
	if _, err := os.Stat(settingsPath); os.IsNotExist(err) {
		needSetup = true
	}

	// check DB file presence
	dbPath, err := utils.GetRiftDBFilePath()
	if err != nil {
		return err
	}
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		needSetup = true
	}

	// check whether the binary is newer than what was last set up
	if !needSetup {
		if settings.RIFTSETTINGS == nil || IsVersionGreater(constant.APPVERSION, settings.RIFTSETTINGS.Version) {
			needSetup = true
		}
	}

	if needSetup {
		message := style.RenderStringWithColor(i18n.LANGUAGEMAPPING.RiftAutoSetupTriggered, style.ColorCyanSoft, false)
		logger.LOGGER.LogToTerminal([]string{message, ""})
		if err := RiftSetup(); err != nil {
			return err
		}
		// stamp the new version so we don't re-run setup on the next invocation
		settings.UpdateVersion(constant.APPVERSION)
	}

	return nil
}

// ----------------------------------
//
//	Reports whether binaryVersion is a valid semver string that is strictly
//	greater than settingsVersion. Returns false if either string is not valid semver.
//
// ----------------------------------
func IsVersionGreater(binaryVersion, settingsVersion string) bool {
	return semver.IsValid(binaryVersion) && semver.IsValid(settingsVersion) && semver.Compare(binaryVersion, settingsVersion) > 0
}

// ----------------------------------
//
//	The list of keywords reserved by rift that cannot be used as waypoint names.
//	These mirror rift's own subcommands and future planned commands.
//
// ----------------------------------
var ReservedCommandKeywords = []string{
	constant.RIFT_CMD_KEYWORD,
	constant.AWAKEN_CMD_KEYWORD,
	constant.DISCOVER_CMD_KEYWORD,
	constant.WAYPOINT_CMD_KEYWORD,
	constant.SPELL_CMD_KEYWORD,
	constant.LEARN_CMD_KEYWORD,
	constant.SPELLBOOK_CMD_KEYWORD,
	constant.CAST_CMD_KEYWORD,
	constant.RITUAL_CMD_KEYWORD,
	constant.INSCRIBE_CMD_KEYWORD,
	constant.SCROLL_CMD_KEYWORD,
	constant.SORCERY_CMD_KEYWORD,
	constant.SUMMON_CMD_KEYWORD,
	constant.DEPLOY_CMD_KEYWORD,
	constant.RUNE_CMD_KEYWORD,
	constant.SEER_CMD_KEYWORD,
	constant.RECALL_CMD_KEYWORD,
	constant.LOOT_CMD_KEYWORD,
	constant.GRIMOIRE_CMD_KEYWORD,
	constant.LORE_CMD_KEYWORD,
	constant.STATS_CMD_KEYWORD,
}

// ----------------------------------
//
//	Returns an error if the given waypoint name conflicts with a reserved rift
//	keyword (e.g. "awaken", "discover"). This prevents waypoints from shadowing
//	rift's own subcommands at the shell level.
//
// ----------------------------------
func CheckIfKeywordIsReservedForRift(arg string) error {
	if slices.Contains(ReservedCommandKeywords, arg) {
		errorMessage := style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.RiftReservedKeywordError, arg), style.ColorError, false)
		return fmt.Errorf("%s", errorMessage)
	}
	return nil
}

// ----------------------------------
//
//	Retrieves the string value of the named flag from cmd, wraps any error
//	into a user-facing i18n message, and trims surrounding whitespace from
//	the returned value.
//
// ----------------------------------
func GetFlagString(cmd *cobra.Command, flagName string) (string, error) {
	value, err := cmd.Flags().GetString(flagName)
	if err != nil {
		return "", fmt.Errorf("%s", style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.RiftFlagRetrieveError, flagName, err.Error()), style.ColorError, false))
	}
	return strings.TrimSpace(value), nil
}

// ----------------------------------
//
//	Returns an error if the given nickname contains whitespace. Nicknames must
//	be a single token with no embedded whitespace so they can be used as shell
//	identifiers without quoting.
//
// ----------------------------------
func IsNickNameValid(nickname string) error {
	if strings.IndexFunc(nickname, unicode.IsSpace) != -1 || utf8.RuneCountInString(nickname) < 1 {
		return fmt.Errorf("%s", style.RenderStringWithColor(i18n.LANGUAGEMAPPING.InvalidValueProvided, style.ColorError, false))
	}

	return nil
}
