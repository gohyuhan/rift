package awaken

import (
	"github.com/gohyuhan/rift/api/utils"
	"github.com/spf13/cobra"
)

// ----------------------------------
//
//	Cobra handler for the awaken command. Triggers the full rift setup flow.
//
// ----------------------------------
var RiftAwakenFunc = func(cmd *cobra.Command, args []string) error {
	return utils.RiftSetup()
}
