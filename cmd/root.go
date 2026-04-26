package cmd

import (
	"fmt"
	"os"

	"github.com/gohyuhan/rift/api/root"
	"github.com/gohyuhan/rift/constant"
	"github.com/gohyuhan/rift/i18n"
	"github.com/spf13/cobra"
)

var rootKeyword = fmt.Sprintf("%s [waypoint name]", constant.RIFT_CMD_KEYWORD)

var rootCmd = &cobra.Command{
	Use:  rootKeyword,
	Args: cobra.MaximumNArgs(1),
	RunE: root.RiftRootFunc,
}

func init() {
	// Redirect all cobra output (help, usage, errors) to stderr so the shell
	// wrapper never tries to eval anything other than an intentional cd command.
	rootCmd.SetOut(os.Stderr)
	rootCmd.SetErr(os.Stderr)

	// settings related flag
	rootCmd.Flags().String("language", "", "")
	rootCmd.Flags().Bool("autoupdate", false, "")
	rootCmd.Flags().Bool("download-pre-release", false, "")
	rootCmd.Flags().Bool("update", false, "")
	rootCmd.Flags().Bool("version", false, "")

	// feature flag
	rootCmd.Flags().String("cast", "", "")
}

// ----------------------------------
//
//	Sets the short descriptions for the root and all subcommands from the
//	active i18n mapping. Must be called before Execute.
//
// ----------------------------------
func InitCmdI18n() {
	rootCmd.Short = i18n.LANGUAGEMAPPING.RiftDescription
	rootCmd.Long = constant.APPLOGO
	rootCmd.Flags().Lookup("language").Usage = i18n.LANGUAGEMAPPING.RiftFlagLanguageDescription
	rootCmd.Flags().Lookup("autoupdate").Usage = i18n.LANGUAGEMAPPING.RiftFlagAutoUpdateDescription
	rootCmd.Flags().Lookup("download-pre-release").Usage = i18n.LANGUAGEMAPPING.RiftFlagDownloadPreReleaseDescription
	rootCmd.Flags().Lookup("update").Usage = i18n.LANGUAGEMAPPING.RiftFlagUpdateDescription
	rootCmd.Flags().Lookup("version").Usage = i18n.LANGUAGEMAPPING.RiftFlagVersionDescription
	rootCmd.Flags().Lookup("cast").Usage = i18n.LANGUAGEMAPPING.RiftFlagCastDescription
	initAwakenI18n()
	initDiscoverI18n()
	initWaypointI18n()
	initLearnI18n()
	initSpellI18n()
	initSpellbookI18n()
	initRuneI18n()
	initInscribeI18n()
	initRitualI18n()
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
