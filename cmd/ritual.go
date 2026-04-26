package cmd

import (
	"fmt"

	"github.com/gohyuhan/rift/api/ritual"
	"github.com/gohyuhan/rift/constant"
	"github.com/gohyuhan/rift/i18n"
	"github.com/spf13/cobra"
)

var ritualKeyword = fmt.Sprintf("%s [ritual name]", constant.RITUAL_CMD_KEYWORD)

var ritualCmd = &cobra.Command{
	Use:  ritualKeyword,
	Args: cobra.ExactArgs(1),
	RunE: ritual.RiftRitualFunc,
}

// ----------------------------------
//
//	Registers the ritual subcommand under the root command.
//
// ----------------------------------
func init() {
	ritualCmd.Flags().Bool("forget", false, "")
	rootCmd.AddCommand(ritualCmd)
}

// ----------------------------------
//
//	Sets the ritual command's short description from the active i18n mapping.
//
// ----------------------------------
func initRitualI18n() {
	ritualCmd.Short = i18n.LANGUAGEMAPPING.RiftRitualDescription
	ritualCmd.Flags().Lookup("forget").Usage = i18n.LANGUAGEMAPPING.RiftFlagRitualForgetDescription
}
