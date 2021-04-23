package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newDoctorCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "doctor",
		Short: "Print diagnostics and check the system for potential problems",
		Long:  `doctor checks the system for potential problems and prints environmental stuff useful for debugging purposes. Exits with a non-zero status if any potential problems are found.`,
		Args:  cobra.MaximumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(fmt.Sprintf("dapla-cli %v", versionInfo()))
			fmt.Println("\nConfig:")
			fmt.Println(effectiveConfig())
			fmt.Println("\nAPIs:")
			fmt.Println(allAPIUrlsString())
		},
	}
}

func init() {
	doctorCommand := newDoctorCommand()
	rootCmd.AddCommand(doctorCommand)
}
