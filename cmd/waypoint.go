package cmd

import (
	"github.com/gohyuhan/rift/api/waypoint"
	"github.com/gohyuhan/rift/i18n"
	"github.com/spf13/cobra"
)

const waypointKeyword = "waypoint [waypoint name <optional>]"

var waypointCmd = &cobra.Command{
	Use:  waypointKeyword,
	Args: cobra.MaximumNArgs(1),
	RunE: waypoint.RiftWaypointFunc,
}

// ----------------------------------
//
//	Registers the waypoint subcommand under the root command.
//
// ----------------------------------
func init() {
	waypointCmd.Flags().String("rebind", "", "")
	waypointCmd.Flags().String("reforge", "", "")
	waypointCmd.Flags().Bool("destroy", false, "")
	waypointCmd.MarkFlagsMutuallyExclusive("destroy", "rebind", "reforge")

	// if rebind is not follow by arguments, the default is " ", which will be "" after space trim.
	waypointCmd.Flags().Lookup("rebind").NoOptDefVal = " "
	rootCmd.AddCommand(waypointCmd)
}

// ----------------------------------
//
//	Sets the waypoint command's short description from the active i18n mapping.
//
// ----------------------------------
func initWaypointI18n() {
	waypointCmd.Short = i18n.LANGUAGEMAPPING.RiftWaypointDescription
	waypointCmd.Flags().Lookup("destroy").Usage = i18n.LANGUAGEMAPPING.RiftFlagWaypointDestroyDescription
	waypointCmd.Flags().Lookup("rebind").Usage = i18n.LANGUAGEMAPPING.RiftFlagWaypointRebindDescription
	waypointCmd.Flags().Lookup("reforge").Usage = i18n.LANGUAGEMAPPING.RiftFlagWaypointReforgeDescription
}
