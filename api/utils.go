package api

import (
	"fmt"
	"os"
	"slices"

	"github.com/gohyuhan/rift/constant"
	"github.com/gohyuhan/rift/db"
	"github.com/gohyuhan/rift/i18n"
	"github.com/gohyuhan/rift/internal/shell"
	"github.com/gohyuhan/rift/logger"
	"github.com/gohyuhan/rift/settings"
	"github.com/gohyuhan/rift/style"
	"github.com/gohyuhan/rift/utils"
	"golang.org/x/mod/semver"
)

// ----------------------------------
//
//	Runs the full rift setup: checks binary is in PATH, detects the shell,
//	installs the shell integration if not already present, and initializes the DB.
//
// ----------------------------------
func RiftSetup() error {
	if !shell.BinaryInPath() {
		message := style.RenderStringWithColor(i18n.LANGUAGEMAPPING.BinaryNotInPath, style.ColorYellowWarm, false)
		logger.LOGGER.LogToTerminal([]string{message})
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
		message := style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.ShellAlreadyInstalled, cfgFile), style.ColorGreenSoft, false)
		logger.LOGGER.LogToTerminal([]string{message})
	} else {
		if err := shell.Install(sh); err != nil {
			return err
		}
	}

	// just to change to new line
	logger.LOGGER.LogToTerminal([]string{""})

	dbSetupErr := db.SetupDB()
	if dbSetupErr != nil {
		return dbSetupErr
	}

	return nil
}

// ----------------------------------
//
//	CheckAndRunSetup checks whether settings.json exists, the db exists,
//	and whether the binary version is newer than the recorded settings version.
//	If any condition is true, RiftSetup is run and the settings version is updated.
//
// ----------------------------------
func CheckAndRunSetup() error {
	needSetup := false

	settingsPath, err := utils.GetRiftSettingsFilePath()
	if err != nil {
		return err
	}
	if _, err := os.Stat(settingsPath); os.IsNotExist(err) {
		needSetup = true
	}

	dbPath, err := utils.GetRiftDBFilePath()
	if err != nil {
		return err
	}
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		needSetup = true
	}

	if !needSetup {
		if settings.RIFTSETTINGS == nil || isVersionGreater(constant.APPVERSION, settings.RIFTSETTINGS.Version) {
			needSetup = true
		}
	}

	if needSetup {
		message := style.RenderStringWithColor(i18n.LANGUAGEMAPPING.RiftAutoSetupTriggered, style.ColorCyanSoft, false)
		logger.LOGGER.LogToTerminal([]string{message, ""})
		if err := RiftSetup(); err != nil {
			return err
		}
		settings.UpdateVersion(constant.APPVERSION)
	}

	return nil
}

// ----------------------------------
//
//	Returns true if binaryVersion is a valid semver string greater than settingsVersion.
//
// ----------------------------------
func isVersionGreater(binaryVersion, settingsVersion string) bool {
	return semver.IsValid(binaryVersion) && semver.IsValid(settingsVersion) && semver.Compare(binaryVersion, settingsVersion) > 0
}

// ----------------------------------
//
//	The list of keywords reserved by rift that cannot be used as waypoint names.
//
// ----------------------------------
var ReservedCommandKeywords = []string{
	"rift",
	"awaken",
	"discover",
	"waypoint",
	"spell",
	"spellbook",
	"cast",
	"ritual",
	"scroll",
	"sorcery",
	"summon",
	"deploy",
	"rune",
	"seer",
	"recall",
	"loot",
	"grimore",
	"lore",
	"stats",
}

// ----------------------------------
//
// This is to check if those waypoint name defined by the user didn't conflict with rift's reserved keyword,
// such as `awaken`.
//
// ----------------------------------
func CheckIfKeywordIsReservedForRift(arg string) error {
	if slices.Contains(ReservedCommandKeywords, arg) {
		errorMessage := style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.RiftReservedKeywordError, arg), style.ColorError, false)
		return fmt.Errorf("%s", errorMessage)
	}
	return nil
}
