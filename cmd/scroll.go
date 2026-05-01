package cmd

import (
	"fmt"

	"github.com/gohyuhan/rift/api/scroll"
	"github.com/gohyuhan/rift/constant"
	"github.com/gohyuhan/rift/i18n"
	"github.com/spf13/cobra"
)

var scrollKeyword = fmt.Sprintf("%s [ritual name <optional>]", constant.SCROLL_CMD_KEYWORD)

var scrollCmd = &cobra.Command{
	Use:  scrollKeyword,
	Args: cobra.MaximumNArgs(1),
	RunE: scroll.RiftScrollFunc,
}

// ----------------------------------
//
//	Registers the scroll subcommand under the root command.
//
// ----------------------------------
func init() {
	rootCmd.AddCommand(scrollCmd)
}

// ----------------------------------
//
//	Sets the scroll command's short description from the active i18n mapping.
//
// ----------------------------------
func initScrollI18n() {
	scrollCmd.Short = i18n.LANGUAGEMAPPING.RiftScrollDescription
}
