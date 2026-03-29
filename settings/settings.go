package settings

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"slices"
	"strings"
	"time"

	"github.com/gohyuhan/rift/constant"
	"github.com/gohyuhan/rift/i18n"
	"github.com/gohyuhan/rift/logger"
	"github.com/gohyuhan/rift/style"
	"github.com/gohyuhan/rift/utils"
)

var RIFTSETTINGS *RiftSettings

type RiftSettings struct {
	Version             string    `json:"version"`
	LanguageCode        string    `json:"language_code"`
	LastUpdateCheckTime time.Time `json:"last_update_check_time"`
	AutoUpdate          bool      `json:"auto_update"`
	DownloadPreRelease  bool      `json:"download_pre_release"`
}

var RiftDefaultConfigSettings = RiftSettings{
	Version:             constant.APPVERSION,
	LanguageCode:        "EN",
	LastUpdateCheckTime: time.Now().UTC(),
	AutoUpdate:          true,
	DownloadPreRelease:  false,
}

// ----------------------------------
//
//	Loads the settings file from disk into RIFTSETTINGS and initializes the i18n
//	mapping. If the file does not exist, a default settings file is written first.
//	If the file is unreadable or unparseable, a default settings file is written
//	and used instead. Any fields that are zero-valued or invalid are repaired and
//	persisted before the settings object is stored in RIFTSETTINGS.
//
// ----------------------------------
func InitOrReadSettings() {
	RIFTSETTINGS = &RiftDefaultConfigSettings
	i18n.InitRiftLanguageMapping(RIFTSETTINGS.LanguageCode)

	settingsPath, err := utils.GetRiftSettingsFilePath()
	if err != nil {
		return
	}

	// If config doesn't exist, create a default one
	if _, err := os.Stat(settingsPath); os.IsNotExist(err) {
		writeDefaultSettings(settingsPath)
		return
	}

	data, err := os.ReadFile(settingsPath)
	if err != nil {
		errorMessage := style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.SettingsReadError, err.Error()), style.ColorError, false)
		logger.LOGGER.LogToTerminal([]string{errorMessage})
		writeDefaultSettings(settingsPath)
		return
	}

	var settings RiftSettings
	if err := json.Unmarshal(data, &settings); err != nil {
		errorMessage := style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.SettingsParseError, err.Error()), style.ColorError, false)
		logger.LOGGER.LogToTerminal([]string{errorMessage})
		writeDefaultSettings(settingsPath)
		return
	}

	// Validate and fix missing or invalid fields
	changed := ensureConfigIntegrity(&settings, &RiftDefaultConfigSettings)
	if changed {
		saveSettings(settingsPath, settings)
	}

	RIFTSETTINGS = &settings
	i18n.InitRiftLanguageMapping(RIFTSETTINGS.LanguageCode)
}

// ----------------------------------
//
//	ensureConfigIntegrity checks every field against the default.
//	If a field is zero or invalid (type mismatch), it assigns the default value.
//
// ----------------------------------
func ensureConfigIntegrity(cfg *RiftSettings, def *RiftSettings) bool {
	cfgVal := reflect.ValueOf(cfg).Elem()
	defVal := reflect.ValueOf(def).Elem()
	changed := false

	for i := 0; i < cfgVal.NumField(); i++ {
		field := cfgVal.Field(i)
		defaultField := defVal.Field(i)

		switch field.Kind() {
		case reflect.String:
			if field.String() == "" {
				field.SetString(defaultField.String())
				changed = true
			}
		case reflect.Int, reflect.Int64:
			if field.Int() == 0 {
				field.SetInt(defaultField.Int())
				changed = true
			}
		case reflect.Float64:
			if field.Float() == 0 {
				field.SetFloat(defaultField.Float())
				changed = true
			}
		default:
			// for unsupported types, just reset if zero
			if reflect.DeepEqual(field.Interface(), reflect.Zero(field.Type()).Interface()) {
				field.Set(defaultField)
				changed = true
			}
		}
	}
	return changed
}

// ----------------------------------
//
//	Sets LastUpdateCheckTime to the current UTC time and persists the change
//	to disk. Errors from the settings path lookup are silently ignored.
//
// ----------------------------------
func UpdateLastFetchTime() {
	RIFTSETTINGS.LastUpdateCheckTime = time.Now().UTC()
	settingsPath, err := utils.GetRiftSettingsFilePath()
	if err == nil {
		saveSettings(settingsPath, *RIFTSETTINGS)
	}
}

// ----------------------------------
//
//	Serializes settings to indented JSON and writes it to settingsPath,
//	overwriting any existing file. Logs a user-visible error on failure.
//
// ----------------------------------
func saveSettings(settingsPath string, settings RiftSettings) {
	file, err := os.Create(settingsPath)
	if err != nil {
		errorMessage := style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.SettingsWriteError, err.Error()), style.ColorError, false)
		logger.LOGGER.LogToTerminal([]string{errorMessage})
		return
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")
	_ = enc.Encode(settings)
}

// ----------------------------------
//
//	Writes RiftDefaultConfigSettings to settingsPath via saveSettings.
//
// ----------------------------------
func writeDefaultSettings(settingsPath string) {
	saveSettings(settingsPath, RiftDefaultConfigSettings)
}

// ----------------------------------
//
//	Validates languageCode against the supported list, then persists it and
//	reinitializes the i18n mapping. Logs an error if the code is unsupported
//	or the settings path cannot be resolved.
//
// ----------------------------------
func UpdateLanguageCode(languageCode string) {
	langCode := strings.ToUpper(languageCode)
	if !slices.Contains(i18n.SUPPORTED_LANGUAGE_CODE, langCode) {
		errorMessage := style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.SettingsLanguageNotSupported, langCode, strings.Join(i18n.SUPPORTED_LANGUAGE_CODE, ", ")), style.ColorError, false)
		logger.LOGGER.LogToTerminal([]string{errorMessage})
		return
	}
	RIFTSETTINGS.LanguageCode = langCode
	settingsPath, err := utils.GetRiftSettingsFilePath()
	if err == nil {
		saveSettings(settingsPath, *RIFTSETTINGS)
		// update the current language mapping
		i18n.InitRiftLanguageMapping(RIFTSETTINGS.LanguageCode)
		message := style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.SettingsLanguageUpdated, langCode), style.ColorGreenSoft, false)
		logger.LOGGER.LogToTerminal([]string{message})
	} else {
		errorMessage := style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.SettingsPathError, err.Error()), style.ColorError, false)
		logger.LOGGER.LogToTerminal([]string{errorMessage})
	}
}

// ----------------------------------
//
//	Updates the Version field in RIFTSETTINGS and persists the change to disk.
//	Errors from the settings path lookup are silently ignored.
//
// ----------------------------------
func UpdateVersion(version string) {
	RIFTSETTINGS.Version = version
	settingsPath, err := utils.GetRiftSettingsFilePath()
	if err == nil {
		saveSettings(settingsPath, *RIFTSETTINGS)
	}
}

// ----------------------------------
//
//	Updates the DownloadPreRelease field in RIFTSETTINGS, persists the change,
//	and logs a confirmation. Logs an error if the settings path cannot be resolved.
//
// ----------------------------------
func UpdateDownloadPreRelease(downloadPreRelease bool) {
	RIFTSETTINGS.DownloadPreRelease = downloadPreRelease
	settingsPath, err := utils.GetRiftSettingsFilePath()
	if err == nil {
		saveSettings(settingsPath, *RIFTSETTINGS)
		message := style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.SettingsDownloadPreReleaseUpdated, downloadPreRelease), style.ColorGreenSoft, false)
		logger.LOGGER.LogToTerminal([]string{message})
	} else {
		errorMessage := style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.SettingsPathError, err.Error()), style.ColorError, false)
		logger.LOGGER.LogToTerminal([]string{errorMessage})
	}
}

// ----------------------------------
//
//	Updates the AutoUpdate field in RIFTSETTINGS, persists the change, and
//	logs a confirmation. Logs an error if the settings path cannot be resolved.
//
// ----------------------------------
func UpdateAutoUpdate(autoUpdate bool) {
	RIFTSETTINGS.AutoUpdate = autoUpdate
	settingsPath, err := utils.GetRiftSettingsFilePath()
	if err == nil {
		saveSettings(settingsPath, *RIFTSETTINGS)
		message := style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.SettingsAutoUpdateUpdated, autoUpdate), style.ColorGreenSoft, false)
		logger.LOGGER.LogToTerminal([]string{message})
	} else {
		errorMessage := style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.SettingsPathError, err.Error()), style.ColorError, false)
		logger.LOGGER.LogToTerminal([]string{errorMessage})
	}
}
