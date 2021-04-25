package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/statisticsnorway/dapla-cli/export"
)

var (
	req           export.Request
	pseudoRuleMap map[string]string
)

func newExportCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "export [PATH]",
		Short: "Export a dataset",
		Long:  `The export command exports (and optionally depseudonymizes) a specified dataset`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			req.DatasetPath = args[0]

			req.PseudoRules = make([]export.PseudoRule, 0, len(pseudoRuleMap))
			i := 0
			for p, f := range pseudoRuleMap {
				i++
				req.PseudoRules = append(req.PseudoRules, export.PseudoRule{
					Name:    fmt.Sprintf("rule-%d", i),
					Pattern: p,
					Func:    f})
			}

			spinner := newSpinner("This might take some time...")
			client := export.NewClient(apiURLOf(APINamePseudoSvc), authToken(), viper.GetBool(CFGDebug))
			res, err := client.Export(req)
			cobra.CheckErr(err)
			spinner.Stop()

			fmt.Println(res.TargetURI)
		},
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return doAutoComplete(toComplete)
		},
	}
}

func init() {
	exportCommand := newExportCommand()
	exportCommand.Flags().StringVarP(&req.TargetContentName, "name", "n", "", "optional descriptive name of the contents, used as baseline for the target archive name")
	exportCommand.Flags().StringVarP(&req.TargetPassword, "password", "p", "", "password used to protect target archive")
	exportCommand.MarkFlagRequired("password")
	exportCommand.Flags().BoolVar(&req.Depseudonymize, "depseudo", false, "depseudonymize data during export")
	exportCommand.Flags().StringToStringVar(&pseudoRuleMap, "pseudo-rules", map[string]string{}, "explicit pseudo rules to use")
	exportCommand.Flags().StringVar(&req.PseudoRulesDatasetPath, "pseudo-rules-path", "", "path to retrieve pseudo rules from")

	// TODO: Add validation rule that fails if both pseudo-rules and pseudo-rules-path flags are specified

	rootCmd.AddCommand(exportCommand)
}
