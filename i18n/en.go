package i18n

var eN = LanguageMapping{
	// General
	ConfigPathError:          "Invalid configuration path, [ERROR: %s]",
	RiftReservedKeywordError: "`%s` is a reserved keyword for rift",
	RiftDetectedShell:        "rift: detected shell: %s",
	CWDIsNotDirError:         "Current working directory is not a valid directory",
	PathNotAbsoluteError:     "Path must be absolute, got: %s",
	NotFileOrDirError:        "Path does not exist as a file or directory",

	// Settings related
	SettingsPathError:                 "Failed to access settings directory, [ERROR: %s]",
	SettingsReadError:                 "Failed to read settings file, [ERROR: %s]",
	SettingsParseError:                "Failed to parse settings file, resetting to defaults, [ERROR: %s]",
	SettingsWriteError:                "Failed to write settings file, [ERROR: %s]",
	SettingsLanguageUpdated:           "rift: language set to %s",
	SettingsLanguageNotSupported:      "rift: language %q is not supported (supported: %s)",
	SettingsAutoUpdateUpdated:         "rift: auto-update set to %t",
	SettingsDownloadPreReleaseUpdated: "rift: download pre-release set to %t",

	// DB related
	DBPathError:                 "Failed to access database directory, [ERROR: %s]",
	DBSetupError:                "Failed to initialize database, [ERROR: %s]",
	DBOpenError:                 "Failed to open database, perhaps you haven't setup rift yet, run `rift awaken` for initialization",
	SettingsBucketNotFoundError: "Settings bucket not found in database, perhaps rerun `rift awaken`",
	WaypointBucketNotFoundError: "Waypoint bucket not found in database, perhaps rerun `rift awaken`",
	WaypointDataCorruptedError:  "Waypoint data for %q is corrupted and could not be read",

	// Updater related
	UpdaterDownloadPrompt:               "A new version %s is available. Download now? (y/n): ",
	UpdaterFailToCheckForUpdate:         "Failed to check for updates: %v",
	UpdaterFailToCreateRequest:          "failed to create request: %v",
	UpdaterFailToFetchRelease:           "failed to fetch latest release: %v",
	UpdaterNoReleasesFound:              "no releases found",
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
	RiftDescription:                       "Navigate path easily by your predefined waypoint name",
	RiftAwakenDescription:                 "Awaken rift within your shell [setup and initialize for the first time for using rift]",
	RiftDiscoverDescription:               "Assign a waypoint name to the current working directory",
	RiftWaypointDescription:               "List all waypoints or display info for a specific waypoint",
	RiftFlagLanguageDescription:           "Set the language for rift (supported: EN, JA, ZH-HANS, ZH-HANT)",
	RiftFlagAutoUpdateDescription:         "Set whether rift should automatically check for updates (use --autoupdate to set true, --autoupdate=false to set false)",
	RiftFlagDownloadPreReleaseDescription: "Set whether rift should also download pre-release versions, or only stable releases (use --download-pre-release to set true, --download-pre-release=false to set false)",
	RiftFlagWaypointDestroyDescription:    "Remove a waypoint by name",
	RiftFlagWaypointRebindDescription:     "Reassign an existing waypoint to a new path; defaults to the current working directory if no path is provided, or uses the given absolute path if valid",
	RiftFlagWaypointReforgeDescription:    "Rename an existing waypoint to a new name",

	// Waypoint related
	RiftSavedWaypoint:                "rift: saved %q -> %s",
	RiftUnknownWaypoint:              "rift: unknown waypoint name %q",
	RiftWaypointAlreadyExistsError:   "Waypoint %q already exists, pointing to %s",
	RiftWaypointDoNotExistsError:     "rift: waypoint %q does not exist",
	RiftWaypointUpdateError:          "rift: failed to update waypoint %q",
	RiftWaypointSealedError:          "rift: waypoint %q is sealed and cannot be travelled to due to %q",
	RiftWaypointSealedLabel:          "(SEALED)",
	RiftWaypointRetrieveAllError:     "rift: failed to retrieve waypoints",
	RiftWaypointDetailName:           "Waypoint Name:",
	RiftWaypointDetailPath:           "Waypoint Path:",
	RiftWaypointDetailDiscovered:     "Waypoint Discovered:",
	RiftWaypointDetailTravelledCount: "Waypoint Travelled Count:",
	RiftWaypointDetailSealed:         "Waypoint Sealed:",
	RiftWaypointDetailSealedReason:   "Sealed Reason:",
	RiftWaypointDetailSealedTrue:     "Yes",
	RiftWaypointDetailSealedFalse:    "No",

	// Flag related
	RiftFlagRetrieveError: "rift: failed to retrieve flag %q, [ERROR: %s]",

	// Setup related
	CheckAndRunSetupError:  "rift: setup failed, [ERROR: %s]",
	RiftAutoSetupTriggered: "rift: settings and config setup triggered automatically",
}
