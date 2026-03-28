package cmd

import (
	"fmt"

	"github.com/gohyuhan/rift/api/discover"
	"github.com/gohyuhan/rift/constant"
	"github.com/gohyuhan/rift/i18n"
	"github.com/spf13/cobra"
)

var discoverKeyword = fmt.Sprintf("%s [waypoint name]", constant.DISCOVER_CMD_KEYWORD)

var discoverCmd = &cobra.Command{
	Use:  discoverKeyword,
	Args: cobra.ExactArgs(1),
	RunE: discover.RiftDiscoverFunc,
}

// ----------------------------------
//
//	Registers the discover subcommand under the root command.
//
// ----------------------------------
func init() {
	rootCmd.AddCommand(discoverCmd)
}

// ----------------------------------
//
//	Sets the discover command's short description from the active i18n mapping.
//
// ----------------------------------
func initDiscoverI18n() {
	discoverCmd.Short = i18n.LANGUAGEMAPPING.RiftDiscoverDescription
}
