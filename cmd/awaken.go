package cmd

import (
	"fmt"
	"os"

	"github.com/gohyuhan/rift/internal/shell"
	"github.com/spf13/cobra"
)

var awakenCmd = &cobra.Command{
	Use:   "awaken",
	Short: "Awaken rift within your shell",
	RunE: func(cmd *cobra.Command, args []string) error {
		if !shell.BinaryInPath() {
			fmt.Fprintln(os.Stderr, "rift: warning: rift binary not found in PATH — make sure it is installed before awakening")
		}

		sh := shell.Detect()
		fmt.Fprintf(os.Stderr, "rift: detected shell: %s\n", sh)
		return shell.Install(sh)
	},
}

func init() {
	rootCmd.AddCommand(awakenCmd)
}
