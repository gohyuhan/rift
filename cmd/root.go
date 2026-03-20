package cmd

import (
	"os"

	"github.com/gohyuhan/rift/api"
	"github.com/gohyuhan/rift/i18n"
	"github.com/spf13/cobra"
)

const rootKeyword = "rift [waypoint name]"

var rootCmd = &cobra.Command{
	Use:  rootKeyword,
	Args: cobra.MaximumNArgs(1),
	RunE: api.RiftRootFunc,
}

func init() {
	// Redirect all cobra output (help, usage, errors) to stderr so the shell
	// wrapper never tries to eval anything other than an intentional cd command.
	rootCmd.SetOut(os.Stderr)
	rootCmd.SetErr(os.Stderr)
}

// ----------------------------------
//
//	Sets the short descriptions for the root and all subcommands from the
//	active i18n mapping. Must be called before Execute.
//
// ----------------------------------
func InitCmdI18n() {
	rootCmd.Short = i18n.LANGUAGEMAPPING.RiftDescription
	initAwakenI18n()
	initDiscoverI18n()
}

// ----------------------------------
//
//	Bootstraps i18n descriptions and runs the root cobra command.
//	Exits with code 1 on error.
//
// ----------------------------------
func Execute() {
	InitCmdI18n()
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
