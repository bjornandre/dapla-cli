package cmd

import (
	"dapla-cli/rest"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"strings"
)

func init() {
	rootCmd.AddCommand(lsCommand)
}

var lsCommand = &cobra.Command{
	Use:   "ls [PATH]...",
	Short: "List information about the dataset(s) under PATH",
	Long:  `TODO`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		if jupyter && bearerToken != "" {
			panic(errors.New("cannot use both --jupyter and --token"))
		}

		var client *rest.Client
		if jupyter {
			client = rest.NewClientWithJupyter(serverUrl)
		}

		if bearerToken != "" {
			client = rest.NewClient(serverUrl, bearerToken)
		}
		if client == nil {
			panic(errors.New("use --jupyter or define the --token"))
		}

		for _, path := range args {
			res, err := client.ListDatasets(path)
			if err != nil {
				panic(err)
			}
			fmt.Println(res)
		}

		fmt.Println("LIST! " + strings.Join(args, " "))
	},
}
