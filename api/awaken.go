package api

import (
	"github.com/spf13/cobra"
)

var RiftAwakenFunc = func(cmd *cobra.Command, args []string) error {
	return RiftSetup()
}
