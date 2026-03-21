package i18n

// -------------------------------------------------------
//
//	Language Data Structure
//	* the sequence and structure will follow EN's
//
// -------------------------------------------------------
type LanguageMapping struct {
	// General
	ConfigPathError          string
	RiftReservedKeywordError string
	RiftDetectedShell        string
	CWDIsNotDirError         string
	PathNotAbsoluteError     string
	NotFileOrDirError        string

	// Settings related
	SettingsPathError                 string
	SettingsReadError                 string
	SettingsParseError                string
	SettingsWriteError                string
	SettingsLanguageUpdated           string
	SettingsLanguageNotSupported      string
	SettingsAutoUpdateUpdated         string
	SettingsDownloadPreReleaseUpdated string

	// DB related
	DBPathError                 string
	DBSetupError                string
	DBOpenError                 string
	SettingsBucketNotFoundError string
	WaypointBucketNotFoundError string
	WaypointDataCorruptedError  string

	// Updater related
	UpdaterDownloadPrompt               string
	UpdaterFailToCheckForUpdate         string
	UpdaterFailToCreateRequest          string
	UpdaterFailToFetchRelease           string
	UpdaterNoReleasesFound              string
	UpdaterFailToReadResponse           string
	UpdaterFailToParseJSON              string
	UpdaterFailToExtractBinary          string
	UpdaterUnsupportedArchiveFormat     string
	UpdaterBinaryNotFoundInArchive      string
	UpdaterAlreadyLatest                string
	UpdaterDownloading                  string
	UpdaterUnSupportedOS                string
	UpdaterDownloadFail                 string
	UpdaterBinaryReplaceFail            string
	UpdaterDownloadSuccess              string
	UpdaterDownloadUnexpectedStatusCode string
	UpdaterRequiresSudo                 string

	// Shell related
	ShellCMDNotSupported  string
	ShellUnsupported      string
	ShellNoConfigFile     string
	ShellAlreadyInstalled string
	ShellInstallSuccess   string
	ShellInstallReload    string
	BinaryNotInPath       string

	// cmd description
	RiftDescription                       string
	RiftAwakenDescription                 string
	RiftDiscoverDescription               string
	RiftWaypointDescription               string
	RiftFlagLanguageDescription           string
	RiftFlagAutoUpdateDescription         string
	RiftFlagDownloadPreReleaseDescription string
	RiftFlagWaypointDestroyDescription    string
	RiftFlagWaypointRebindDescription     string
	RiftFlagWaypointReforgeDescription    string
	RiftFlagUpdateDescription             string

	// Waypoint related
	RiftSavedWaypoint                     string
	RiftUnknownWaypoint                   string
	RiftWaypointAlreadyExistsError        string
	RiftWaypointDoNotExistsError          string
	RiftWaypointUpdateError               string
	RiftWaypointSealedError               string
	RiftWaypointSealedLabel               string
	RiftWaypointRetrieveAllError          string
	RiftWaypointDetailName                string
	RiftWaypointDetailPath                string
	RiftWaypointDetailDiscovered          string
	RiftWaypointDetailTravelledCount      string
	RiftWaypointDetailSealed              string
	RiftWaypointDetailSealedReason        string
	RiftWaypointDetailSealedTrue          string
	RiftWaypointDetailSealedFalse         string
	RiftWaypointDestroySuccess            string
	RiftWaypointDestroyError              string
	RiftWaypointRebindNotDirError         string
	RiftWaypointRebindSuccess             string
	RiftWaypointRebindError               string
	RiftWaypointReforgeEmptyError         string
	RiftWaypointReforgeAlreadyExistsError string
	RiftWaypointReforgeError              string
	RiftWaypointReforgeSuccess            string

	// Flag related
	RiftFlagRetrieveError string

	// Setup related
	CheckAndRunSetupError  string
	RiftAutoSetupTriggered string
}
