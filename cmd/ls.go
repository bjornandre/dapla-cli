package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/statisticsnorway/dapla-cli/rest"
	"io"
	"os"
	"text/tabwriter"
	"time"
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

		// Use newline when not in terminal (piped)
		var printFunction func(datasets *rest.DatasetResponse, output io.Writer)
		if fileInfo, _ := os.Stdout.Stat(); (fileInfo.Mode() & os.ModeCharDevice) != 0 {
			printFunction = printTabular
		} else {
			printFunction = printNewLine
		}

		for _, path := range args {
			res, err := client.ListDatasets(path)
			if err != nil {
				panic(err)
			}
			printFunction(res, os.Stdout)
		}
	},
}

// Prints the dataset names
func printNewLine(datasets *rest.DatasetResponse, output io.Writer) {
	writer := bufio.NewWriter(output)
	defer writer.Flush()
	for _, dataset := range *datasets {
		fmt.Fprintln(writer, dataset.Name)
	}
}

// Prints the datasets in tabular format.
func printTabular(datasets *rest.DatasetResponse, output io.Writer) {
	writer := tabwriter.NewWriter(output, 32, 0, 2, ' ', tabwriter.TabIndent)
	defer writer.Flush()
	fmt.Fprintln(writer, "Name\tAuthor\tCreated")
	for _, dataset := range *datasets {
		fmt.Fprintf(writer, "%s\t%s\t%s\n", dataset.Name, dataset.CreatedBy, dataset.CreatedAt.Format(time.RFC3339))
	}
}
