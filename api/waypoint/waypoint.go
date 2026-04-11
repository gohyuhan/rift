package waypoint

import (
	"fmt"
	"strings"

	apiUtils "github.com/gohyuhan/rift/api/utils"
	"github.com/gohyuhan/rift/api/waypoint/features"
	waypointUI "github.com/gohyuhan/rift/api/waypoint/ui"
	"github.com/gohyuhan/rift/db"
	"github.com/gohyuhan/rift/logger"
	"github.com/spf13/cobra"
)

// ----------------------------------
//
//	Cobra handler for the waypoint command.
//	With no args, launches the interactive TUI; if the user selects a waypoint,
//	prints a cd command for the shell wrapper to eval and increments the travel count.
//	With a waypoint name arg and:
//	  --destroy  : permanently removes the named waypoint from the DB
//	  --rebind   : reassigns the waypoint to a new path (defaults to CWD)
//	  --reforge  : renames the waypoint to a new name
//	  no flag    : shows detailed info for the named waypoint
//
// ----------------------------------
var RiftWaypointFunc = func(cmd *cobra.Command, args []string) error {
	// open DB so we can read waypoint records
	bboltReadDb, bboltReadDbErr := db.OpenReadDB()
	if bboltReadDbErr != nil {
		return bboltReadDbErr
	}
	defer db.CloseDB(bboltReadDb)

	// no args — list all start waypoint interactive UI
	if len(args) < 1 {
		pathToNavigate, waypointName, interactiveErr := waypointUI.RunWaypointInteractive(bboltReadDb)
		if interactiveErr != nil {
			return interactiveErr
		}

		if pathToNavigate != "" && waypointName != "" {
			// Only this line goes to stdout — the shell wrapper evals it.
			fmt.Printf("cd %q", pathToNavigate)

			// best-effort: increment travel count; failure is silently ignored
			apiUtils.UpdateWaypointTravelledCount(waypointName)
		}

		return nil
	}

	// extract and normalize the waypoint name from the first argument
	waypointName := strings.TrimSpace(args[0])

	// check which mutually exclusive operation flag was provided
	destroyFlagCalled := cmd.Flags().Changed("destroy")
	rebindFlagCalled := cmd.Flags().Changed("rebind")
	reforgeFlagCalled := cmd.Flags().Changed("reforge")

	if destroyFlagCalled {
		// if destroy flag is called, we destroy the discovered waypoint in the waypoint bucket
		destroyWaypointErr := features.DestroyDiscoveredWaypoint(waypointName, true)

		if destroyWaypointErr != nil {
			return destroyWaypointErr
		}
	} else if rebindFlagCalled {
		// retrieve the target path from the flag value; empty string signals "use CWD"
		rebindTo, rebindToErr := apiUtils.GetFlagString(cmd, "rebind")
		if rebindToErr != nil {
			return rebindToErr
		}
		rebindWaypointErr := features.RebindWaypoint(waypointName, rebindTo, true)

		if rebindWaypointErr != nil {
			return rebindWaypointErr
		}
	} else if reforgeFlagCalled {
		// retrieve the new waypoint name from the flag value
		reforgeTo, reforgeToErr := apiUtils.GetFlagString(cmd, "reforge")
		if reforgeToErr != nil {
			return reforgeToErr
		}
		reforgeWaypointErr := features.ReforgeWaypoint(waypointName, reforgeTo, true)

		if reforgeWaypointErr != nil {
			return reforgeWaypointErr
		}
	} else {
		retrieveWaypointInfoDetail, retrieveWaypointInfoDetailErr := features.RetrieveWaypointInfoDetail(bboltReadDb, waypointName)

		if retrieveWaypointInfoDetailErr != nil {
			return retrieveWaypointInfoDetailErr
		}

		logger.LOGGER.LogToTerminal(retrieveWaypointInfoDetail)
	}

	return nil
}
