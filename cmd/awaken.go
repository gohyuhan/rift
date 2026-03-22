package cmd

import (
	"github.com/gohyuhan/rift/api/awaken"
	"github.com/gohyuhan/rift/i18n"
	"github.com/spf13/cobra"
)

const awakenKeyword = "awaken"

var awakenCmd = &cobra.Command{
	Use:  awakenKeyword,
	RunE: awaken.RiftAwakenFunc,
}

// ----------------------------------
//
//	Registers the awaken subcommand under the root command.
//
// ----------------------------------
func init() {
	rootCmd.AddCommand(awakenCmd)
}

// ----------------------------------
//
//	Sets the awaken command's short description from the active i18n mapping.
//
// ----------------------------------
func initAwakenI18n() {
	awakenCmd.Short = i18n.LANGUAGEMAPPING.RiftAwakenDescription
}
