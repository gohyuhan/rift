package cmd

import (
	"fmt"

	"github.com/gohyuhan/rift/api/spell"
	"github.com/gohyuhan/rift/constant"
	"github.com/gohyuhan/rift/i18n"
	"github.com/spf13/cobra"
)

var spellKeyword = fmt.Sprintf("%s [spell name]", constant.SPELL_CMD_KEYWORD)

var spellCmd = &cobra.Command{
	Use:  spellKeyword,
	Args: cobra.ExactArgs(1),
	RunE: spell.RiftSpellFunc,
}

// ----------------------------------
//
//	Registers the spell subcommand under the root command.
//
// ----------------------------------
func init() {
	spellCmd.Flags().Bool("forget", false, "")
	rootCmd.AddCommand(spellCmd)
}

// ----------------------------------
//
//	Sets the Spell command's short description from the active i18n mapping.
//
// ----------------------------------
func initSpellI18n() {
	spellCmd.Short = i18n.LANGUAGEMAPPING.RiftSpellDescription
	spellCmd.Flags().Lookup("forget").Usage = i18n.LANGUAGEMAPPING.RiftFlagSpellForgetDescription
}
