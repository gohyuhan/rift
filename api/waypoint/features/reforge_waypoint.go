package features

import (
	"fmt"

	apiUtils "github.com/gohyuhan/rift/api/utils"
	"github.com/gohyuhan/rift/i18n"
	"github.com/gohyuhan/rift/logger"
	"github.com/gohyuhan/rift/style"
	"go.etcd.io/bbolt"
)

// ----------------------------------
//
//	Renames the named waypoint to a new name.
//	If reforgeTo is empty, returns an error immediately — a waypoint name is
//	required and there is no sensible default.
//	If a waypoint with the target name already exists, the operation is
//	rejected to prevent silently overwriting an existing record.
//	The rename is atomic within the Update transaction: the record is written
//	under the new key first, then the old key is deleted. All fields (path,
//	sealed state, travelled count, timestamps) are preserved unchanged.
//
// ----------------------------------
func ReforgeWaypoint(bboltDb *bbolt.DB, waypointName string, reforgeTo string, logToTerminal bool) error {
	// validate that a non-empty new name was provided before opening the Update transaction;
	// unlike rebind, there is no sensible default — an empty name is always an error
	if reforgeTo == "" {
		return fmt.Errorf("%s", style.RenderStringWithColor(i18n.LANGUAGEMAPPING.RiftWaypointReforgeEmptyError, style.ColorError, false))
	}

	reforgeErr := bboltDb.Update(func(tx *bbolt.Tx) error {
		// fetch the current waypoint record and its bucket in a single helper call
		waypointBucket, waypoint, retrieveErr := apiUtils.GetWaypointForUpdate(tx, waypointName)
		if retrieveErr != nil {
			return retrieveErr
		}

		// check if a waypoint with the new name already exists to prevent overwriting
		existingWaypoint := waypointBucket.Get([]byte(reforgeTo))
		if existingWaypoint != nil {
			return fmt.Errorf("%s", style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.RiftWaypointReforgeAlreadyExistsError, reforgeTo), style.ColorError, false))
		}

		// write the record under the new name first; if this fails the old key is untouched
		waypoint.WaypointName = reforgeTo
		putWaypointErr := apiUtils.PutWaypoint(waypointBucket, reforgeTo, waypoint)
		if putWaypointErr != nil {
			return putWaypointErr
		}

		// remove the old key only after the new one is confirmed written
		destroyWaypointErr := waypointBucket.Delete([]byte(waypointName))
		if destroyWaypointErr != nil {
			return fmt.Errorf("%s", style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.RiftWaypointReforgeError, waypointName, destroyWaypointErr.Error()), style.ColorError, false))
		}

		// report the rename to the terminal
		if logToTerminal {
			message := style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.RiftWaypointReforgeSuccess, waypointName, reforgeTo), style.ColorGreenSoft, false)
			logger.LOGGER.LogToTerminal([]string{message})
		}
		return nil
	})

	return reforgeErr
}
