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
var (
	serverUrl   string
	bearerToken string
	jupyter     bool
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&serverUrl, "server", "s", "",
		"set URI of the API server")
	rootCmd.PersistentFlags().StringVar(&bearerToken, "token", "",
		"set the Bearer token to use to authenticate with the server")
	rootCmd.PersistentFlags().BoolVar(&jupyter, "jupyter", false,
		"fetch the Bearer token from jupyter")
}
