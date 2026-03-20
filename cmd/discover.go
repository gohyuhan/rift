package cmd

import (
	"github.com/gohyuhan/rift/api"
	"github.com/gohyuhan/rift/i18n"
	"github.com/spf13/cobra"
)

const discoverKeyword = "discover [waypoint name]"

var discoverCmd = &cobra.Command{
	Use:  discoverKeyword,
	Args: cobra.MaximumNArgs(1),
	RunE: api.RiftDiscoverFunc,
}

func init() {
	rootCmd.AddCommand(discoverCmd)
}

func initDiscoverI18n() {
	discoverCmd.Short = i18n.LANGUAGEMAPPING.RiftDiscoverDescription
}
