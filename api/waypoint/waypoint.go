package waypoint

import (
	"strings"

	apiUtils "github.com/gohyuhan/rift/api/utils"
	"github.com/gohyuhan/rift/db"
	"github.com/gohyuhan/rift/logger"
	"github.com/spf13/cobra"
)

// ----------------------------------
//
//	Cobra handler for the waypoint command.
//	With no args, lists every stored waypoint (name + path, sealed state).
//	With a waypoint name arg and:
//	  --destroy  : permanently removes the named waypoint from the DB
//	  --rebind   : reassigns the waypoint to a new path (defaults to CWD)
//	  --reforge  : renames the waypoint to a new name
//	  no flag    : shows detailed info for the named waypoint
//
// ----------------------------------
var RiftWaypointFunc = func(cmd *cobra.Command, args []string) error {
	// open DB so we can read waypoint records
	bboltDB, bboltDBErr := db.OpenDB()
	if bboltDBErr != nil {
		return bboltDBErr
	}
	defer db.CloseDB(bboltDB)

	// no args — list all waypoints
	if len(args) < 1 {
		allWaypointInfo, allWaypointInfoErr := waypointInteractive(bboltDB)
		if allWaypointInfoErr != nil {
			return allWaypointInfoErr
		}

		logger.LOGGER.LogToTerminal(allWaypointInfo)

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
		destroyWaypointErr := destroyDiscoveredWaypoint(bboltDB, waypointName)

		if destroyWaypointErr != nil {
			return destroyWaypointErr
		}
	} else if rebindFlagCalled {
		// retrieve the target path from the flag value; empty string signals "use CWD"
		rebindTo, rebindToErr := apiUtils.GetFlagString(cmd, "rebind")
		if rebindToErr != nil {
			return rebindToErr
		}
		rebindWaypointErr := rebindWaypoint(bboltDB, waypointName, rebindTo)

		if rebindWaypointErr != nil {
			return rebindWaypointErr
		}
	} else if reforgeFlagCalled {
		// retrieve the new waypoint name from the flag value
		reforgeTo, reforgeToErr := apiUtils.GetFlagString(cmd, "reforge")
		if reforgeToErr != nil {
			return reforgeToErr
		}
		reforgeWaypointErr := reforgeWaypoint(bboltDB, waypointName, reforgeTo)

		if reforgeWaypointErr != nil {
			return reforgeWaypointErr
		}
	} else {
		retrieveWaypointInfoDetail, retrieveWaypointInfoDetailErr := retrieveWaypointInfoDetail(bboltDB, waypointName)

		if retrieveWaypointInfoDetailErr != nil {
			return retrieveWaypointInfoDetailErr
		}

		logger.LOGGER.LogToTerminal(retrieveWaypointInfoDetail)
	}

	return nil
}
