package cmd

import (
	"fmt"

	"github.com/gohyuhan/rift/api/spellbook"
	"github.com/gohyuhan/rift/constant"
	"github.com/gohyuhan/rift/i18n"
	"github.com/spf13/cobra"
)

var spellbookKeyword = fmt.Sprintf("%s [spell name <optional>]", constant.SPELLBOOK_CMD_KEYWORD)

var spellbookCmd = &cobra.Command{
	Use:  spellbookKeyword,
	Args: cobra.MaximumNArgs(1),
	RunE: spellbook.RiftSpellbookFunc,
}

// ----------------------------------
//
//	Registers the spellbook subcommand under the root command.
//
// ----------------------------------
func init() {
	rootCmd.AddCommand(spellbookCmd)
}

// ----------------------------------
//
//	Sets the spellbook command's short description from the active i18n mapping.
//
// ----------------------------------
func initSpellbookI18n() {
	spellbookCmd.Short = i18n.LANGUAGEMAPPING.RiftSpellbookDescription
}
