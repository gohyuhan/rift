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

	// Settings related
	SettingsPathError  string
	SettingsReadError  string
	SettingsParseError string
	SettingsWriteError string

	// DB related
	DBPathError                   string
	DBSetupError                  string
	DBOpenError                   string
	SettingsBucketNotFoundError   string
	CheckpointBucketNotFoundError string

	// Updater related
	UpdaterDownloadPrompt               string
	UpdaterFailToCheckForUpdate         string
	UpdaterFailToCreateRequest          string
	UpdaterFailToFetchRelease           string
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
	RiftDescription       string
	RiftAwakenDescription string

	// cmd root
	RiftSavedCheckpoint   string
	RiftUnknownCheckpoint string
}
