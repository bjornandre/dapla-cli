package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dapla",
	Short: "dapla command line utility",
	Long: `The dapla command is a collection of utilities you can use with the dapla
				platform.`,
}

func Execute() error {
	return rootCmd.Execute()
}
