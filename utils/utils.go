package utils

import (
	"fmt"
	"os"
	"path/filepath"

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

func GetCWD() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return cwd, nil
}

func CheckIsDir(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	if info.IsDir() {
		return true, nil
	}
	return false, nil
}

func CheckIsPathExist(path string) (bool, error) {
	_, pathErr := os.Stat(path)
	if pathErr == nil {
		return true, nil
	}
	if os.IsNotExist(pathErr) {
		return false, nil
	}
	return false, pathErr
}
