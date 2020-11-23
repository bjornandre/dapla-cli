package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/statisticsnorway/dapla-cli/rest"
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
		switch {
		case jupyter:
			var err error
			client, err = rest.NewClientWithJupyter(serverUrl)
			if err != nil {
				panic(err)
			}
		case bearerToken != "":
			client = rest.NewClient(serverUrl, bearerToken)
		default:
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
