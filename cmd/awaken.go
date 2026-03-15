package cmd

import (
	"github.com/gohyuhan/rift/api"
	"github.com/gohyuhan/rift/i18n"
	"github.com/spf13/cobra"
)

const awakenKeyword = "awaken"

var awakenCmd = &cobra.Command{
	Use:  awakenKeyword,
	RunE: api.RiftAwakenFunc,
}

func init() {
	rootCmd.AddCommand(awakenCmd)
}

func initAwakenI18n() {
	awakenCmd.Short = i18n.LANGUAGEMAPPING.RiftAwakenDescription
}
