package api

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/gohyuhan/rift/constant"
	"github.com/gohyuhan/rift/db"
	"github.com/gohyuhan/rift/i18n"
	"github.com/gohyuhan/rift/internal/shell"
	"github.com/gohyuhan/rift/logger"
	pb "github.com/gohyuhan/rift/proto"
	"github.com/gohyuhan/rift/settings"
	"github.com/gohyuhan/rift/style"
	"github.com/gohyuhan/rift/utils"
	"github.com/spf13/cobra"
	"go.etcd.io/bbolt"
	"golang.org/x/mod/semver"
	"google.golang.org/protobuf/proto"
)

// ----------------------------------
//
//	Runs the full rift setup sequence:
//	  1. Warns if the rift binary is not on PATH
//	  2. Detects the current shell
//	  3. Installs the shell integration into the shell config file (if not already present)
//	  4. Initializes (or migrates) the bbolt DB
//
// ----------------------------------
func RiftSetup() error {
	// warn if the binary won't be found after this session ends
	if !shell.BinaryInPath() {
		message := style.RenderStringWithColor(i18n.LANGUAGEMAPPING.BinaryNotInPath, style.ColorYellowWarm, false)
		logger.LOGGER.LogToTerminal([]string{message})
	}

	// detect the running shell so we know which config file to modify
	sh := shell.Detect()
	msg := style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.RiftDetectedShell, sh), style.ColorPurpleVibrant, false)
	logger.LOGGER.LogToTerminal([]string{msg})

	cfgFile, cfgFileErr := shell.ConfigFile(sh)
	if cfgFileErr != nil {
		// CMD or unsupported shell — Install returns the descriptive error
		return shell.Install(sh)
	}

	// check whether the shell integration snippet is already present
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

	// initialize the DB (creates buckets if they don't exist yet)
	dbSetupErr := db.SetupDB()
	if dbSetupErr != nil {
		return dbSetupErr
	}

	return nil
}

// ----------------------------------
//
//	Decides whether setup needs to run and triggers it if so.
//	Setup is required when any of the following is true:
//	  - the settings file does not exist on disk
//	  - the DB file does not exist on disk
//	  - the binary version is newer than the version recorded in settings
//	After a successful setup the settings version is stamped with the current binary version.
//
// ----------------------------------
func CheckAndRunSetup() error {
	needSetup := false

	// check settings file presence
	settingsPath, err := utils.GetRiftSettingsFilePath()
	if err != nil {
		return err
	}
	if _, err := os.Stat(settingsPath); os.IsNotExist(err) {
		needSetup = true
	}

	// check DB file presence
	dbPath, err := utils.GetRiftDBFilePath()
	if err != nil {
		return err
	}
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		needSetup = true
	}

	// check whether the binary is newer than what was last set up
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
		// stamp the new version so we don't re-run setup on the next invocation
		settings.UpdateVersion(constant.APPVERSION)
	}

	return nil
}

// ----------------------------------
//
//	Reports whether binaryVersion is a valid semver string that is strictly
//	greater than settingsVersion. Returns false if either string is not valid semver.
//
// ----------------------------------
func isVersionGreater(binaryVersion, settingsVersion string) bool {
	return semver.IsValid(binaryVersion) && semver.IsValid(settingsVersion) && semver.Compare(binaryVersion, settingsVersion) > 0
}

// ----------------------------------
//
//	The list of keywords reserved by rift that cannot be used as waypoint names.
//	These mirror rift's own subcommands and future planned commands.
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
//	Returns an error if the given waypoint name conflicts with a reserved rift
//	keyword (e.g. "awaken", "discover"). This prevents waypoints from shadowing
//	rift's own subcommands at the shell level.
//
// ----------------------------------
func CheckIfKeywordIsReservedForRift(arg string) error {
	if slices.Contains(ReservedCommandKeywords, arg) {
		errorMessage := style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.RiftReservedKeywordError, arg), style.ColorError, false)
		return fmt.Errorf("%s", errorMessage)
	}
	return nil
}

// ----------------------------------
//
//	Persists the waypoint name into the corrupted-records bucket via its own
//	Update transaction (so the write commits regardless of what the caller's
//	transaction did), then returns a user-facing corruption error.
//	The Update's own error is intentionally not propagated — recording is
//	best-effort; the corruption error is always returned to the caller.
//
// ----------------------------------
func recordCorruptedWaypointInfo(bboltDB *bbolt.DB, waypointName string) error {
	// best-effort write — ignore the Update error; the caller always gets the corruption message
	bboltDB.Update(func(tx *bbolt.Tx) error {
		waypointCorruptedBucket := tx.Bucket(db.WaypointDataCorruptedBucketRecord)
		if waypointCorruptedBucket != nil {
			waypointCorruptedBucket.Put([]byte(waypointName), []byte(waypointName))
		}
		return nil
	})
	return fmt.Errorf("%s", style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.WaypointDataCorruptedError, waypointName), style.ColorError, false))
}

// ----------------------------------
//
//	Fetches and deserializes the named waypoint from the bucket within an
//	already-open Update transaction. Returns the bucket, the deserialized
//	record, or an error if the bucket is missing, the waypoint does not exist,
//	or the stored proto is corrupted. Callers mutate the returned record and
//	re-persist it via bucket.Put.
//
// ----------------------------------
func getWaypointForUpdate(tx *bbolt.Tx, waypointName string) (*bbolt.Bucket, *pb.Waypoint, error) {
	bucket := tx.Bucket(db.WaypointBucket)
	if bucket == nil {
		return nil, nil, fmt.Errorf("%s", style.RenderStringWithColor(i18n.LANGUAGEMAPPING.WaypointBucketNotFoundError, style.ColorError, false))
	}

	existing := bucket.Get([]byte(waypointName))
	if existing == nil {
		return nil, nil, fmt.Errorf("%s", style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.RiftWaypointDoNotExistsError, waypointName), style.ColorError, false))
	}

	waypoint := &pb.Waypoint{}
	if err := proto.Unmarshal(existing, waypoint); err != nil {
		return nil, nil, fmt.Errorf("%s", style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.WaypointDataCorruptedError, waypointName), style.ColorError, false))
	}

	return bucket, waypoint, nil
}

// ----------------------------------
//
//	Persists a mutated waypoint record back into its bucket. Returns an error
//	if marshalling or the bucket write fails.
//
// ----------------------------------
func putWaypoint(bucket *bbolt.Bucket, waypointName string, waypoint *pb.Waypoint) error {
	data, err := proto.Marshal(waypoint)
	if err != nil {
		return fmt.Errorf("%s", style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.RiftWaypointUpdateError, waypointName), style.ColorError, false))
	}
	return bucket.Put([]byte(waypointName), data)
}

// ----------------------------------
//
//	Retrieves the string value of the named flag from cmd, wraps any error
//	into a user-facing i18n message, and trims surrounding whitespace from
//	the returned value.
//
// ----------------------------------
func getFlagString(cmd *cobra.Command, flagName string) (string, error) {
	value, err := cmd.Flags().GetString(flagName)
	if err != nil {
		return "", fmt.Errorf("%s", style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.RiftFlagRetrieveError, flagName, err.Error()), style.ColorError, false))
	}
	return strings.TrimSpace(value), nil
}
