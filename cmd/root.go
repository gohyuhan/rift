package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var memorize string

var rootCmd = &cobra.Command{
	Use:   "rift [nickname]",
	Short: "Navigate to saved paths by nickname",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// --memorize <nickname>: save current directory
		if memorize != "" {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}
			// TODO: store.Set(memorize, cwd)
			_ = cwd
			fmt.Fprintf(os.Stderr, "rift: saved %q -> %s\n", memorize, cwd)
			return nil
		}

		// rift <nickname>: emit cd command for the shell wrapper to eval
		if len(args) == 0 {
			return cmd.Help()
		}

		// TODO: path, err := store.Get(args[0])
		path := ""
		if path == "" {
			fmt.Fprintf(os.Stderr, "rift: unknown nickname %q\n", args[0])
			os.Exit(1)
		}

		// Only this line goes to stdout — the shell wrapper evals it.
		fmt.Printf("cd %q", path)
		return nil
	},
}

func init() {
	rootCmd.Flags().StringVar(&memorize, "memorize", "", "save current directory under this nickname")
	// Redirect all cobra output (help, usage, errors) to stderr so the shell
	// wrapper never tries to eval anything other than an intentional cd command.
	rootCmd.SetOut(os.Stderr)
	rootCmd.SetErr(os.Stderr)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
