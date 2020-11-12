package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"strings"
)

func init() {
	rootCmd.AddCommand(rmCommand)
}

var rmCommand = &cobra.Command{
	Use:   "rm [PATH]...",
	Short: "Remove the dataset(s) under PATH",
	Long:  `TODO`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("RM! " + strings.Join(args, " "))
	},
}
