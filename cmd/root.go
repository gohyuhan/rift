package cmd

import (
	"fmt"
	"os"

	"github.com/gohyuhan/rift/i18n"
	"github.com/gohyuhan/rift/logger"
	"github.com/gohyuhan/rift/style"
	"github.com/spf13/cobra"
)

var memorize string

var rootCmd = &cobra.Command{
	Use:  "rift [checkpoint name]",
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// --memorize <checkpoint name>: save current directory
		if memorize != "" {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}
			// TODO: store.Set(memorize, cwd)
			_ = cwd
			logger.LOGGER.LogToTerminal([]string{style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.RiftSavedCheckpoint, memorize, cwd), style.ColorGreenSoft, false)})
			return nil
		}

		// rift <checkpoint name>: emit cd command for the shell wrapper to eval
		if len(args) == 0 {
			return cmd.Help()
		}

		// TODO: path, err := store.Get(args[0])
		path := ""
		if path == "" {
			logger.LOGGER.LogToTerminal([]string{style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.RiftUnknownCheckpoint, args[0]), style.ColorError, false)})
			os.Exit(1)
		}

		// Only this line goes to stdout — the shell wrapper evals it.
		fmt.Printf("cd %q", path)
		return nil
	},
}

func init() {
	rootCmd.Flags().StringVar(&memorize, "memorize", "", "save current directory under this checkpoint name")
	// Redirect all cobra output (help, usage, errors) to stderr so the shell
	// wrapper never tries to eval anything other than an intentional cd command.
	rootCmd.SetOut(os.Stderr)
	rootCmd.SetErr(os.Stderr)
}

func InitCmdI18n() {
	rootCmd.Short = i18n.LANGUAGEMAPPING.RiftDescription
	initAwakenI18n()
}

func Execute() {
	InitCmdI18n()
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
