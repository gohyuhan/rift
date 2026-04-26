package cmd

import (
	"fmt"

	"github.com/gohyuhan/rift/api/inscribe"
	"github.com/gohyuhan/rift/constant"
	"github.com/gohyuhan/rift/i18n"
	"github.com/spf13/cobra"
)

var inscribeKeyword = fmt.Sprintf("%s [ritual name] [ritual description] [ritual command(s)]", constant.INSCRIBE_CMD_KEYWORD)

var inscribeCmd = &cobra.Command{
	Use:  inscribeKeyword,
	Args: cobra.ExactArgs(3),
	RunE: inscribe.RiftInscribeFunc,
}

// ----------------------------------
//
//	Registers the inscribe subcommand under the root command.
//
// ----------------------------------
func init() {
	inscribeCmd.Flags().Bool("override", false, "")
	rootCmd.AddCommand(inscribeCmd)
}

// ----------------------------------
//
//	Sets the inscribe command's short description from the active i18n mapping.
//
// ----------------------------------
func initInscribeI18n() {
	inscribeCmd.Short = i18n.LANGUAGEMAPPING.RiftInscribeDescription
	inscribeCmd.Flags().Lookup("override").Usage = i18n.LANGUAGEMAPPING.RiftFlagRitualOverrideDescription
}
