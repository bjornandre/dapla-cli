package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(lsCommand)
}

var lsCommand = &cobra.Command{
	Use:   "ls",
	Short: "List information about the dataset",
	Long:  `TODO`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("LIST!")
	},
}
