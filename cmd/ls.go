package cmd

import (
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
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("LIST! " + strings.Join(args, " "))
	},
}
