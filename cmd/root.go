package cmd

import (
	"os"

	"github.com/gohyuhan/rift/i18n"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:  "rift [waypoint name]",
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	// Redirect all cobra output (help, usage, errors) to stderr so the shell
	// wrapper never tries to eval anything other than an intentional cd command.
	rootCmd.SetOut(os.Stderr)
	rootCmd.SetErr(os.Stderr)
}

func InitCmdI18n() {
	rootCmd.Short = i18n.LANGUAGEMAPPING.RiftDescription
	initAwakenI18n()
	initDiscoverI18n()
}

func Execute() {
	InitCmdI18n()
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
