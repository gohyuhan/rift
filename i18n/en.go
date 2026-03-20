package i18n

var eN = LanguageMapping{
	// General
	ConfigPathError:          "Invalid configuration path, [ERROR: %s]",
	RiftReservedKeywordError: "`%s` is a reserved keyword for rift",
	RiftDetectedShell:        "rift: detected shell: %s",

	// Settings related
	SettingsPathError:  "Failed to access settings directory, [ERROR: %s]",
	SettingsReadError:  "Failed to read settings file, [ERROR: %s]",
	SettingsParseError: "Failed to parse settings file, resetting to defaults, [ERROR: %s]",
	SettingsWriteError: "Failed to write settings file, [ERROR: %s]",

	// DB related
	DBPathError:                 "Failed to access database directory, [ERROR: %s]",
	DBSetupError:                "Failed to initialize database, [ERROR: %s]",
	DBOpenError:                 "Failed to open database, perhaps you haven't setup rift yet, run `rift awaken` for initialization",
	SettingsBucketNotFoundError: "Settings bucket not found in database, perhaps rerun `rift awaken`",
	WaypointBucketNotFoundError: "Waypoint bucket not found in database, perhaps rerun `rift awaken`",

	// Updater related
	UpdaterDownloadPrompt:               "A new version %s is available. Download now? (y/n): ",
	UpdaterFailToCheckForUpdate:         "Failed to check for updates: %v",
	UpdaterFailToCreateRequest:          "failed to create request: %v",
	UpdaterFailToFetchRelease:           "failed to fetch latest release: %v",
	UpdaterFailToReadResponse:           "failed to read response body: %v",
	UpdaterFailToParseJSON:              "failed to parse JSON response: %v",
	UpdaterFailToExtractBinary:          "failed to extract binary: %v",
	UpdaterUnsupportedArchiveFormat:     "unsupported archive format",
	UpdaterBinaryNotFoundInArchive:      "binary not found in archive",
	UpdaterAlreadyLatest:                "You are already on the latest version (%s)",
	UpdaterDownloading:                  "Downloading version %s...",
	UpdaterUnSupportedOS:                "Unsupported OS/architecture: %s/%s",
	UpdaterDownloadFail:                 "Failed to download update: %v",
	UpdaterBinaryReplaceFail:            "Failed to replace binary: %v",
	UpdaterDownloadSuccess:              "Successfully updated to version %s",
	UpdaterDownloadUnexpectedStatusCode: "unexpected status code: %d",
	UpdaterRequiresSudo:                 "Permission denied. Retrying with sudo...",

	// Shell related
	ShellCMDNotSupported:  "Windows Command Prompt does not support shell functions.\nPlease use PowerShell, Git Bash, or WSL instead, then re-run `rift awaken`.",
	ShellUnsupported:      "shell %q is not supported by rift.\nSupported shells: bash, zsh, fish, ksh, PowerShell.\nYou can add the integration manually — see docs/shell-integration.md",
	ShellNoConfigFile:     "shell %q does not have a known config file",
	ShellAlreadyInstalled: "rift: shell integration already present in %s",
	ShellInstallSuccess:   "rift: shell integration added to %s",
	ShellInstallReload:    "rift: restart your shell or run: %s",
	BinaryNotInPath:       "rift: binary not found in PATH — add rift to your PATH to use it after this session",

	// cmd description
	RiftDescription:       "Navigate path easily by your predefined waypoint name",
	RiftAwakenDescription: "Awaken rift within your shell [setup and initialize for the first time for using rift]",

	// cmd root
	RiftSavedWaypoint:   "rift: saved %q -> %s",
	RiftUnknownWaypoint: "rift: unknown waypoint name %q",

	// Setup related
	CheckAndRunSetupError: "rift: setup failed, [ERROR: %s]",
}
