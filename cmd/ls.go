package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"github.com/statisticsnorway/dapla-cli/rest"
	"io"
	"os"
	"strings"
	"text/tabwriter"
	"time"
)

var (
	lsLong bool
)

func newLsCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "ls [PATH]...",
		Short: "List the datasets and folders under a PATH",
		Long: `The ls 	command list the datasets and folders under a given PATH.`,
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			// Create and start spinner
			s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
			s.Color("bgWhite", "bold", "black")
			s.Prefix = "Listing... "
			s.Start()

			var client, err = initClient()
			if err != nil {
				panic(err) // TODO don't panic!
			}

			// Use newline when not in terminal (piped)
			var printFunction func(datasets *rest.ListDatasetResponse, output io.Writer)
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
					exitCode := 1
					switch err.(type) {
					case *rest.HttpError:
						exitCode = 0
					default:
					}
					fmt.Println(err.Error() + "\n")
					os.Exit(exitCode)
				} else if res != nil {
					// Strip the common prefix. Note that we are mutating the
					// elements of res and therefore need to use index notation.
					var prefix = strings.TrimSuffix(path, "/") + "/"
					for i := 0; i < len(*res); i++ {
						(*res)[i].Path = strings.TrimPrefix((*res)[i].Path, prefix)
					}
				} else {
					// TODO what to do if no error and response is nil
				}
			}

			// Stop spinner
			// TODO remove spinner when done
			s.Stop()
		},
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {

			// TODO make test(s)!

			return doAutoComplete(toComplete)
		},
	}
}

func initClient() (*rest.Client, error) {
	if jupyter && bearerToken != "" {
		return nil, errors.New("cannot use both --jupyter and --token")
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
	lsCommand := newLsCommand()
	lsCommand.Flags().BoolVarP(&lsLong, "", "l", false, "use a long listing format")
	rootCmd.AddCommand(lsCommand)
}

// Prints the dataset names
func printNewLine(datasets *rest.ListDatasetResponse, output io.Writer) {
	writer := bufio.NewWriter(output)
	defer writer.Flush()
	for _, dataset := range *datasets {
		fmt.Fprintln(writer, dataset.Path)
	}
}

type ColorWriter struct {
	out io.Writer
}

func (c ColorWriter) Write(p []byte) (int, error) {
	// Simply send the result of replace tag. Note that we need
	// to send the size of the buffer.
	_, err := fmt.Fprint(c.out, color.ReplaceTag(string(p)))
	return len(p), err
}

func printTabular(datasets *rest.ListDatasetResponse, output io.Writer) {
	colorOutput := ColorWriter{out: output}
	// TODO: Test with strip escape (\xff[colorstuff]\xff" ).
	writer := tabwriter.NewWriter(colorOutput, 15, 0, 2, ' ', tabwriter.FilterHTML)
	defer writer.Flush()

	n := 0
	maxColumns := 5
	// Print the folders first.
	for _, dataset := range *datasets {

		if dataset.IsFolder() {
			fmt.Fprintf(writer, "<fg=blue;op=bold;>%s</>/\t", dataset.Path)
			n++
			if n%maxColumns == 0 {
				fmt.Fprintln(writer)
			}
		}
	}
	for _, dataset := range *datasets {
		if dataset.IsDataset() {
			fmt.Fprint(writer, dataset.Path, "\t")
			n++
			if n%maxColumns == 0 {
				fmt.Fprintln(writer)
			}
		}
	}
	_, _ = fmt.Fprintln(writer)
}

// Prints the datasets in tabular format. Datasets are white and folders blue and with a trailing '/'
func printTabularDetails(datasets *rest.ListDatasetResponse, output io.Writer) {
	colorOutput := ColorWriter{out: output}
	writer := tabwriter.NewWriter(colorOutput, 15, 0, 2, ' ', tabwriter.FilterHTML)
	defer writer.Flush()

	fmt.Fprintln(writer,
		"<bold>Name</>\t"+
			"<bold>Author</>\t"+
			"<bold>Created</>\t"+
			"<bold>Type</>\t"+
			"<bold>Valuation</>\t"+
			"<bold>State</>\t")
	for _, dataset := range *datasets {
		if dataset.IsFolder() {
			fmt.Fprintln(writer,
				color.WrapTag(dataset.Path, "blue")+"/", "\t",
				dataset.CreatedBy+"\t",
				dataset.CreatedAt.Format(time.RFC3339), "\t",
				dataset.Type, "\t",
				dataset.Valuation, "\t",
				dataset.State, "\t")
		} else {
			fmt.Fprintln(writer,
				dataset.Path, "\t",
				dataset.CreatedBy, "\t",
				dataset.CreatedAt.Format(time.RFC3339), "\t",
				dataset.Type, "\t",
				dataset.Valuation, "\t",
				dataset.State, "\t")
		}
	}
}
