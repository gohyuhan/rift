// ----------------------------------
//
//	GENERAL UTILS
//
// ----------------------------------

package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/gohyuhan/rift/constant"
	"github.com/gohyuhan/rift/i18n"
	"github.com/gohyuhan/rift/logger"
	"github.com/gohyuhan/rift/style"
)

// ----------------------------------
//
//	getRiftConfigDirPath returns the config path (creates directories if needed)
//
//	*Example on MacOs : /Users/<USER_NAME>/Library/Application Support/rift/
//
// ----------------------------------
func getRiftConfigDirPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		errorMessage := style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.ConfigPathError, err.Error()), style.ColorError, false)
		logger.LOGGER.LogToTerminal([]string{errorMessage})
		return "", fmt.Errorf("%s", errorMessage)
	}
	appDir := filepath.Join(dir, constant.APPNAME)
	if err := os.MkdirAll(appDir, 0o755); err != nil {
		errorMessage := style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.ConfigPathError, err.Error()), style.ColorError, false)
		logger.LOGGER.LogToTerminal([]string{errorMessage})
		return "", fmt.Errorf("%s", errorMessage)
	}

	return appDir, nil
}

// ----------------------------------
//
//	Returns the full file path to the rift database, creating the db
//	subdirectory under the config dir if it does not exist.
//
// ----------------------------------
func GetRiftDBFilePath() (string, error) {
	configPath, configPathErr := getRiftConfigDirPath()
	if configPathErr != nil {
		return "", configPathErr
	}
	dbDirPath := filepath.Join(configPath, "db")
	if err := os.MkdirAll(dbDirPath, 0o755); err != nil {
		errorMessage := style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.DBPathError, err.Error()), style.ColorError, false)
		logger.LOGGER.LogToTerminal([]string{errorMessage})
		return "", fmt.Errorf("%s", errorMessage)
	}

	dbPath := filepath.Join(dbDirPath, constant.APPDBNAME)
	return dbPath, nil
}

// ----------------------------------
//
//	Returns the full file path to the rift settings JSON file, creating the
//	settings subdirectory under the config dir if it does not exist.
//
// ----------------------------------
func GetRiftSettingsFilePath() (string, error) {
	configPath, configPathErr := getRiftConfigDirPath()
	if configPathErr != nil {
		return "", configPathErr
	}
	settingsDirPath := filepath.Join(configPath, "settings")
	if err := os.MkdirAll(settingsDirPath, 0o755); err != nil {
		errorMessage := style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.SettingsPathError, err.Error()), style.ColorError, false)
		logger.LOGGER.LogToTerminal([]string{errorMessage})
		return "", fmt.Errorf("%s", errorMessage)
	}

	settingsPath := filepath.Join(settingsDirPath, constant.APPSETTINGSNAME)
	return settingsPath, nil
}

// ----------------------------------
//
//	Returns the current working directory path.
//
// ----------------------------------
func GetCWD() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return cwd, nil
}

// ----------------------------------
//
//	Reports whether the given absolute path points to a directory.
//	Returns an error if the path is not absolute.
//
// ----------------------------------
func CheckIsDir(path string) (bool, error) {
	if !filepath.IsAbs(path) {
		errorMessage := style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.PathNotAbsoluteError, path), style.ColorError, false)
		return false, fmt.Errorf("%s", errorMessage)
	}
	info, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	if info.IsDir() {
		return true, nil
	}
	return false, nil
}

// ----------------------------------
//
//	Reports whether the given absolute path exists on the filesystem.
//	Returns an error if the path is not absolute. Returns false for
//	non-existent paths and an error only for unexpected stat failures.
//
// ----------------------------------
func CheckIsPathExist(path string) (bool, error) {
	if !filepath.IsAbs(path) {
		errorMessage := style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.PathNotAbsoluteError, path), style.ColorError, false)
		return false, fmt.Errorf("%s", errorMessage)
	}
	_, pathErr := os.Stat(path)
	if pathErr == nil {
		return true, nil
	}
	if os.IsNotExist(pathErr) {
		return false, fmt.Errorf("%s", i18n.LANGUAGEMAPPING.NotFileOrDirError)
	}
	return false, pathErr
}

// ----------------------------------
//
//	Reports whether the first token of cmd is a shell built-in command.
//	Returns false if cmd is empty or the first token is empty.
//
// ----------------------------------
func IsShellBuiltInCmd(cmd []string) bool {
	if len(cmd) < 1 || len(cmd[0]) < 1 {
		return false
	}
	return slices.Contains(constant.ShellBuildInCmd, cmd[0])
}

// ----------------------------------
//
//	riftSubcommands is the exhaustive list of rift's own subcommand keywords.
//	A second token matching any of these is a rift subcommand invocation, not
//	a waypoint navigation.
//
// ----------------------------------
var riftSubcommands = []string{
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
//	IsRiftNavigationCommand reports whether cmd represents a rift waypoint
//	navigation invocation of the form `rift <waypointName>`.
//
//	It returns true only when:
//	  - exactly two tokens are present
//	  - the first token is "rift"
//	  - the second token is not one of rift's own root flags
//	  - the second token is not one of rift's own reserved subcommands
//
// ----------------------------------
func IsRiftNavigationCommand(cmd []string) bool {
	if len(cmd) != 2 {
		return false
	}
	if cmd[0] != constant.RIFT_CMD_KEYWORD {
		return false
	}
	arg := cmd[1]

	return !slices.Contains(riftSubcommands, arg) && !slices.Contains(constant.RiftRootFlags, arg)
}
