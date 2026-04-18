package cmd

import (
	"fmt"

	"github.com/gohyuhan/rift/api/rune"
	"github.com/gohyuhan/rift/constant"
	"github.com/gohyuhan/rift/i18n"
	"github.com/spf13/cobra"
)

var runeKeyword = fmt.Sprintf("%s [waypoint name]", constant.RUNE_CMD_KEYWORD)

var runeCmd = &cobra.Command{
	Use:  runeKeyword,
	Args: cobra.ExactArgs(1),
	RunE: rune.RiftRuneFunc,
}

// ----------------------------------
//
//	Registers the rune subcommand under the root command.
//
// ----------------------------------
func init() {
	rootCmd.AddCommand(runeCmd)
}

// ----------------------------------
//
//	Sets the rune command's short description from the active i18n mapping.
//
// ----------------------------------
func initRuneI18n() {
	runeCmd.Short = i18n.LANGUAGEMAPPING.RiftRuneDescription
}
