package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/juju/ansiterm"
	"github.com/spf13/cobra"
	"github.com/statisticsnorway/dapla-cli/rest"
	"io"
	"os"
	"strings"
	"time"
)

var (
	lsLong bool
)

var lsCommand = &cobra.Command{
	Use:   "ls [PATH]...",
	Short: "List information about the dataset(s) under PATH",
	Long:  `TODO`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		var client, err = initClient()
		if err != nil {
			panic(err)
		}

		// Use newline when not in terminal (piped)
		var printFunction func(datasets *rest.DatasetResponse, output io.Writer)
		if fileInfo, _ := os.Stdout.Stat(); (fileInfo.Mode() & os.ModeCharDevice) != 0 {
			if lsLong {
				printFunction = printTabularDetails
			} else {
				printFunction = printTabular
			}
		} else {
			printFunction = printNewLine
		}

		for _, path := range args {
			res, err := client.ListDatasets(path)
			if err != nil {
				if strings.HasSuffix(err.Error(), "404") {
					fmt.Printf("Cannot access %s: No such dataset or folder", path)
				} else {
					panic(err) //TODO don't panic
				}
			} else {
				printFunction(res, os.Stdout)
			}
		}
	},
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {

		// TODO make test(s)!

		var client, err = initClient()
		if err != nil {
			panic(err)
		}
		var res *rest.DatasetResponse

		if toComplete == "/" {
			res, err = client.ListDatasets(toComplete)
			if err != nil {
				return nil, cobra.ShellCompDirectiveError
			}
			var result []string
			for _, element := range *res {
				result = append(result, element.Path)
			}
			return result, cobra.ShellCompDirectiveNoFileComp
		} else if strings.HasPrefix(toComplete, "/") {
			var result []string

			// Ask for list without last element
			var parentPath = toComplete[0:strings.LastIndex(toComplete, "/")]
			res, err = client.ListDatasets(parentPath)
			if err != nil {
				return nil, cobra.ShellCompDirectiveError
			}

			// Check if last element is a valid path / dataset
			var partialPath = toComplete[strings.LastIndex(toComplete, "/")+1:]
			for _, element := range *res {
				// We have a complete match, ask data-maintenance for elements on that path
				if toComplete == element.Path {
					res, err = client.ListDatasets(toComplete)
					if err != nil {
						return nil, cobra.ShellCompDirectiveError
					} else {
						for _, element := range *res {
							result = append(result, element.Path)
						}
						return result, cobra.ShellCompDirectiveNoFileComp
					}
				} else { // find all elements that matches the last element in the provided path
					var lastPart = element.Path[strings.LastIndex(element.Path, "/")+1 : len(element.Path)]
					if strings.HasPrefix(lastPart, partialPath) {
						// TODO add trailing / to folder element
						result = append(result, element.Path)
					}
				}
			}
			return result, cobra.ShellCompDirectiveNoFileComp
		} else {
			return []string{"/"}, cobra.ShellCompDirectiveNoFileComp
		}

	},
}

func initClient() (*rest.Client, error) {
	if jupyter && bearerToken != "" {
		panic(errors.New("cannot use both --jupyter and --token"))
	}

	switch {

	case jupyter:
		return rest.NewClientWithJupyter(serverUrl)

	case bearerToken != "":
		return rest.NewClient(serverUrl, bearerToken), nil
	default:
		return nil, errors.New("use --jupyter or define the --token")
	}
}

func init() {
	lsCommand.Flags().BoolVarP(&lsLong, "", "l", false, "use a long listing format")
	rootCmd.AddCommand(lsCommand)
}

// Prints the dataset names
func printNewLine(datasets *rest.DatasetResponse, output io.Writer) {
	writer := bufio.NewWriter(output)
	defer writer.Flush()
	for _, dataset := range *datasets {
		fmt.Fprintln(writer, dataset.Path)
	}
}

func printTabular(datasets *rest.DatasetResponse, output io.Writer) {
	folderContext := ansiterm.Foreground(ansiterm.Blue)
	folderContext.SetStyle(ansiterm.Bold)
	datasetContext := ansiterm.Foreground(ansiterm.Default)

	writer := ansiterm.NewTabWriter(output, 15, 0, 2, ' ', 0)
	defer writer.Flush()

	// Print the folders first.
	for _, dataset := range *datasets {
		if dataset.Depth > 0 {
			folderContext.Fprintf(writer, "%s", dataset.Path)
			datasetContext.Fprint(writer, "/\t")
		}
	}
	for _, dataset := range *datasets {
		if dataset.Depth == 0 {
			datasetContext.Fprintf(writer, "%s\t", dataset.Path)
		}
	}
}

// Prints the datasets in tabular format. Datasets are white and folders blue and with a trailing '/'
func printTabularDetails(datasets *rest.DatasetResponse, output io.Writer) {
	writer := ansiterm.NewTabWriter(output, 32, 0, 2, ' ', 0)
	headerContext := ansiterm.Foreground(ansiterm.BrightCyan)
	headerContext.SetStyle(ansiterm.Bold)
	datasetContext := ansiterm.Foreground(ansiterm.White)
	folderContext := ansiterm.Foreground(ansiterm.Blue)
	folderContext.SetStyle(ansiterm.Italic)
	defer writer.Flush()
	headerContext.Fprint(writer, "Name\tAuthor\tCreated\tType\tValuation\tState\n")
	for _, dataset := range *datasets {
		if dataset.Depth > 0 { // is folder
			folderContext.Fprintf(writer, "%s\t%s\t%s\t%s\t%s\t%s\n",
				dataset.Path+"/",
				dataset.CreatedBy,
				dataset.CreatedAt.Format(time.RFC3339),
				dataset.Type,
				dataset.Valuation,
				dataset.State)
		} else {
			datasetContext.Fprintf(writer, "%s\t%s\t%s\t%s\t%s\t%s\n",
				dataset.Path,
				dataset.CreatedBy,
				dataset.CreatedAt.Format(time.RFC3339),
				dataset.Type,
				dataset.Valuation,
				dataset.State)
		}
	}
}
