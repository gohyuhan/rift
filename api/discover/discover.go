package discover

import (
	"fmt"
	"strings"

	apiUtils "github.com/gohyuhan/rift/api/utils"
	"github.com/gohyuhan/rift/i18n"
	"github.com/gohyuhan/rift/logger"
	"github.com/gohyuhan/rift/style"
	"github.com/gohyuhan/rift/utils"
	"github.com/spf13/cobra"
)

// ----------------------------------
//
//	Cobra handler for the discover command.
//	Resolves the current working directory, validates it is a real directory,
//	guards against reserved keyword names, then persists the waypoint mapping
//	(name → CWD) into the DB.
//
// ----------------------------------
var RiftDiscoverFunc = func(command *cobra.Command, args []string) error {
	// resolve and validate the current working directory
	cwd, cwdErr := utils.GetCWD()
	if cwdErr != nil {
		return cwdErr
	}
	isDir, isDirErr := utils.CheckIsDir(cwd)
	if isDirErr != nil {
		return isDirErr
	}

	if !isDir {
		errorMessage := style.RenderStringWithColor(i18n.LANGUAGEMAPPING.CWDIsNotDirError, style.ColorError, false)
		return fmt.Errorf("%s", errorMessage)
	}

	waypointName := strings.TrimSpace(args[0])

	// reject names that clash with rift's own subcommands
	if err := apiUtils.CheckIfKeywordIsReservedForRift(waypointName); err != nil {
		return err
	}

	// reject names that contain spaces
	if err := apiUtils.IsNickNameValid(waypointName); err != nil {
		return err
	}

	saveWaypointErr := saveWaypoint(waypointName, cwd)

	if saveWaypointErr != nil {
		return saveWaypointErr
	}

	message := style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.RiftSavedWaypoint, waypointName, cwd), style.ColorGreenSoft, false)
	logger.LOGGER.LogToTerminal([]string{message})

	return nil
}
