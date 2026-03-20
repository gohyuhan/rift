package settings

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
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
//	InitOrReadConfig loads existing config, ensures schema correctness, or creates default.
//
// ----------------------------------
func InitOrReadSettings() {
	RIFTSETTINGS = &RiftDefaultConfigSettings

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
		errMsg := style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.SettingsReadError, err.Error()), style.ColorError, false)
		logger.LOGGER.LogToTerminal([]string{errMsg})
		writeDefaultSettings(settingsPath)
		return
	}

	var sttings RiftSettings
	if err := json.Unmarshal(data, &sttings); err != nil {
		errMsg := style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.SettingsParseError, err.Error()), style.ColorError, false)
		logger.LOGGER.LogToTerminal([]string{errMsg})
		writeDefaultSettings(settingsPath)
		return
	}

	// Validate and fix missing or invalid fields
	changed := ensureConfigIntegrity(&sttings, &RiftDefaultConfigSettings)
	if changed {
		saveSettings(settingsPath, sttings)
	}

	RIFTSETTINGS = &sttings
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
//	Update and persist the last update check time to current UTC time
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
//	Persist the given config settings to disk as JSON
//
// ----------------------------------
func saveSettings(settingsPath string, settings RiftSettings) {
	file, err := os.Create(settingsPath)
	if err != nil {
		errMsg := style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.SettingsWriteError, err.Error()), style.ColorError, false)
		logger.LOGGER.LogToTerminal([]string{errMsg})
		return
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")
	_ = enc.Encode(settings)
}

// ----------------------------------
//
//	Write the default config settings to disk
//
// ----------------------------------
func writeDefaultSettings(settingsPath string) {
	saveSettings(settingsPath, RiftDefaultConfigSettings)
}

// ----------------------------------
//
//	Update and persist the language code setting
//
// ----------------------------------
func UpdateLanguageCode(languageCode string) {
	RIFTSETTINGS.LanguageCode = strings.ToUpper(languageCode)
	settingsPath, err := utils.GetRiftSettingsFilePath()
	if err == nil {
		saveSettings(settingsPath, *RIFTSETTINGS)
	}
}

// ----------------------------------
//
//	Update version in settings record
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
//	Update download pre release setting
//
// ----------------------------------
func UpdateDownloadPreRelease(downloadPreRelease bool) {
	RIFTSETTINGS.DownloadPreRelease = downloadPreRelease
	settingsPath, err := utils.GetRiftSettingsFilePath()
	if err == nil {
		saveSettings(settingsPath, *RIFTSETTINGS)
	}
}
