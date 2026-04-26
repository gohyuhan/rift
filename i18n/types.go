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
	InvalidValueProvided     string
	SkippingDueToExecutorErr string
	ExecCommandError         string

	// Settings
	SettingsPathError                 string
	SettingsReadError                 string
	SettingsParseError                string
	SettingsWriteError                string
	SettingsLanguageUpdated           string
	SettingsLanguageNotSupported      string
	SettingsAutoUpdateUpdated         string
	SettingsDownloadPreReleaseUpdated string

	// Database
	DBPathError                 string
	DBSetupError                string
	DBOpenError                 string
	WaypointBucketNotFoundError string
	WaypointDataCorruptedError  string
	SpellBucketNotFoundError    string
	SpellDataCorruptedError     string
	RuneBucketNotFoundError     string
	RuneDataCorruptedError      string
	RitualBucketNotFoundError   string
	RitualDataCorruptedError    string

	// Updater
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

	// Shell
	ShellCMDNotSupported  string
	ShellUnsupported      string
	ShellNoConfigFile     string
	ShellAlreadyInstalled string
	ShellInstallSuccess   string
	ShellReloadHint       string
	BinaryNotInPath       string

	// Commands and flags
	RiftDescription                       string
	RiftAwakenDescription                 string
	RiftDiscoverDescription               string
	RiftWaypointDescription               string
	RiftLearnDescription                  string
	RiftSpellDescription                  string
	RiftSpellbookDescription              string
	RiftFlagLanguageDescription           string
	RiftFlagAutoUpdateDescription         string
	RiftFlagDownloadPreReleaseDescription string
	RiftFlagWaypointDestroyDescription    string
	RiftFlagSpellForgetDescription        string
	RiftFlagWaypointRebindDescription     string
	RiftFlagWaypointReforgeDescription    string
	RiftFlagUpdateDescription             string
	RiftFlagVersionDescription            string
	RiftFlagCastDescription               string
	RiftFlagRetrieveError                 string
	RiftRuneDescription                   string
	RiftInscribeDescription               string
	RiftFlagRitualOverrideDescription     string
	RiftRitualDescription                 string
	RiftFlagRitualForgetDescription       string

	// Spell operations
	RiftSpellSaved                      string
	RiftSpellUpdated                    string
	RiftSpellForgetSuccess              string
	RiftSpellForgetError                string
	RiftSpellDoNotExistsError           string
	RiftSpellUpdateError                string
	ForbiddenShellBuiltinSpellCommand   string
	ForbiddenRiftNavigationSpellCommand string
	ForbiddenRiftNavigationRuneCommand  string
	SpellCommandEmpty                   string
	InvalidSpellCommandError            string

	// Rune operations
	RiftRuneEngraveSuccessful string
	RiftRuneEngraveNone       string
	RiftRuneUpdateError       string

	// Ritual operations
	RiftRitualSaved                               string
	RiftRitualDoNotExistsError                    string
	RiftRitualUpdateError                         string
	ForbiddenRiftNavigationRitualCommand          string
	RitualCommandsInvalidDueToShellBuildInCommand string
	RitualCommandEmpty                            string
	RitualExistAndOverrideNotAllowedError         string
	RiftRitualForgetSuccess                       string
	RiftRitualForgetError                         string

	// Waypoint operations
	RiftSavedWaypoint                     string
	RiftWaypointAlreadyExistsError        string
	RiftWaypointDoNotExistsError          string
	RiftWaypointUpdateError               string
	RiftWaypointSealedError               string
	RiftWaypointSealedLabel               string
	RiftWaypointRetrieveAllError          string
	RiftWaypointDestroySuccess            string
	RiftWaypointDestroyError              string
	RiftWaypointRebindNotDirError         string
	RiftWaypointRebindSuccess             string
	RiftWaypointReforgeEmptyError         string
	RiftWaypointReforgeAlreadyExistsError string
	RiftWaypointReforgeError              string
	RiftWaypointReforgeSuccess            string

	// Spell detail view
	RiftSpellDetailName      string
	RiftSpellDetailCommand   string
	RiftSpellDetailAddedAt   string
	RiftSpellDetailCastCount string

	// Waypoint detail view
	RiftWaypointDetailName           string
	RiftWaypointDetailPath           string
	RiftWaypointDetailDiscovered     string
	RiftWaypointDetailTravelledCount string
	RiftWaypointDetailSealed         string
	RiftWaypointDetailSealedReason   string
	RiftWaypointDetailSealedTrue     string
	RiftWaypointDetailSealedFalse    string

	// Waypoint interactive UI
	WaypointInfoListTitle                         string
	WaypointInteractiveError                      string
	RebindPathInputPlaceHolder                    string
	WaypointRebindTitle                           string
	ReforgeWaypointNameInputPlaceHolder           string
	WaypointReforgeTitle                          string
	WaypointUIUpKeyHelp                           string
	WaypointUIUpKeyHelpDescription                string
	WaypointUIDownKeyHelp                         string
	WaypointUIDownKeyHelpDescription              string
	WaypointUIQuitKeyHelp                         string
	WaypointUIQuitKeyHelpDescription              string
	WaypointUIHelpKeyHelp                         string
	WaypointUIHelpKeyHelpDescription              string
	WaypointNavigateKeyHelp                       string
	WaypointNavigateKeyHelpDescription            string
	WaypointDestroyKeyHelp                        string
	WaypointDestroyKeyHelpDescription             string
	WaypointUnsealKeyHelp                         string
	WaypointUnsealKeyHelpDescription              string
	WaypointRebindKeyHelp                         string
	WaypointRebindKeyHelpDescription              string
	WaypointReforgeKeyHelp                        string
	WaypointReforgeKeyHelpDescription             string
	WaypointNameCopyPathCopyKeyHelp               string
	WaypointNameCopyPathCopyKeyHelpDescription    string
	WaypointCopyFromInputValueKeyHelp             string
	WaypointCopyFromInputValueKeyHelpDescription  string
	WaypointPasteIntoInputValueKeyHelp            string
	WaypointPasteIntoInputValueKeyHelpDescription string
	WaypointClosePopUp                            string
	WaypointClosePopUpDescription                 string

	// Spell interactive UI
	SpellInfoListTitle                       string
	SpellbookInteractiveError                string
	RiftSpellRetrieveAllError                string
	SpellUIUpKeyHelp                         string
	SpellUIUpKeyHelpDescription              string
	SpellUIDownKeyHelp                       string
	SpellUIDownKeyHelpDescription            string
	SpellUIQuitKeyHelp                       string
	SpellUIQuitKeyHelpDescription            string
	SpellUIHelpKeyHelp                       string
	SpellUIHelpKeyHelpDescription            string
	SpellCastKeyHelp                         string
	SpellCastKeyHelpDescription              string
	SpellLearnKeyHelp                        string
	SpellLearnKeyHelpDescription             string
	SpellForgetKey                           string
	SpellForgetKeyDescription                string
	SpellClosePopUp                          string
	SpellClosePopUpDescription               string
	SpellUILearnKeyHelp                      string
	SpellUINextInputKeyHelp                  string
	SpellUIPreviousInputKeyHelp              string
	SpellNameInputPlaceHolder                string
	SpellCommandInputPlaceHolder             string
	SpellNameInputTitle                      string
	SpellCommandInputTitle                   string
	SpellUIChooseCastLocationKeyHelp         string
	SpellUIChooseWaypointCastLocationKeyHelp string

	// Rune interactive UI
	RuneInteractiveError                        string
	RuneEngraveTypeOptionListTitle              string
	EngraveRuneEnterTitle                       string
	EngraveRuneLeaveTitle                       string
	EngraveRuneEngraveButton                    string
	RuneCommandsPlaceHolder                     string
	RuneCommandsInvalidDueToShellBuildInCommand string
	EngraveRuneEnterOptionName                  string
	EngraveRuneEnterOptionDesc                  string
	EngraveRuneLeaveOptionName                  string
	EngraveRuneLeaveOptionDesc                  string
	RemoveRuneEnterOptionName                   string
	RemoveRuneEnterOptionDesc                   string
	RemoveRuneLeaveOptionName                   string
	RemoveRuneLeaveOptionDesc                   string

	// Cast location option popup
	CastLocationOptionTitle               string
	CastLocationOptionCurrent             string
	CastLocationOptionCurrentDescription  string
	CastLocationOptionWaypoint            string
	CastLocationOptionWaypointDescription string
	CastWaypointLocationOptionTitle       string

	// Setup
	CheckAndRunSetupError  string
	RiftAutoSetupTriggered string
}
