package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(rmCommand)
}

var rmCommand = &cobra.Command{
	Use:   "rm",
	Short: "Remove dataset",
	Long:  `TODO`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("RM!")
	},
}
