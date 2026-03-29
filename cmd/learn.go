package cmd

import (
	"fmt"

	"github.com/gohyuhan/rift/api/learn"
	"github.com/gohyuhan/rift/constant"
	"github.com/gohyuhan/rift/i18n"
	"github.com/spf13/cobra"
)

var learnKeyword = fmt.Sprintf("%s [spell name] [spell command]", constant.LEARN_CMD_KEYWORD)

var learnCmd = &cobra.Command{
	Use:  learnKeyword,
	Args: cobra.ExactArgs(2),
	RunE: learn.RiftLearnFunc,
}

// ----------------------------------
//
//	Registers the learn subcommand under the root command.
//
// ----------------------------------
func init() {
	rootCmd.AddCommand(learnCmd)
}

// ----------------------------------
//
//	Sets the learn command's short description from the active i18n mapping.
//
// ----------------------------------
func initLearnI18n() {
	learnCmd.Short = i18n.LANGUAGEMAPPING.RiftLearnDescription
}
