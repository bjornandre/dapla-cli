package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/statisticsnorway/dapla-cli/maintenance"
)

var (
	rmDryRun    bool
	rmRecursive bool
)

func init() {
	rmCommand := newRmCommand()
	rmCommand.Flags().BoolVarP(&rmRecursive, "recursive", "", false, "delete recursively")
	rmCommand.Flags().BoolVarP(&rmDryRun, "dry-run", "", false, "dry run")
	rootCmd.AddCommand(rmCommand)
}

func newRmCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "rm [PATH]...",
		Short: "Delete dataset(s)",
		Long:  `The rm command deletes all the version of a given dataset`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			path := args[0]
			if rmRecursive {
				deleteRecursively(path, cmd.OutOrStdout(), cmd.InOrStdin(), rmDryRun)
			} else {
				doDelete(path, cmd.OutOrStdout(), rmDryRun)
			}
		},
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {

			return doAutoComplete(toComplete)

		},
	}
}

func doDelete(path string, output io.Writer, dryRun bool) {
	// Create and start spinner
	spinner := newSpinner("Deleting dataset " + path)
	var client = maintenance.NewClient(apiURLOf(APINameDataMaintenanceSvc), authToken())
	res, err := client.DeleteDatasets(path, dryRun)
	spinner.Stop()

	if err != nil {
		exitCode := 1
		switch err.(type) {
		case *maintenance.HTTPError:
			exitCode = 0
		default:
		}
		fmt.Fprintln(output, err.Error()+"\n")
		os.Exit(exitCode)
	} else {
		printDeleteResponse(res, output, dryRun)
	}
}

func deleteRecursively(path string, output io.Writer, input io.Reader, dryRun bool) {
	var client = maintenance.NewClient(apiURLOf(APINameDataMaintenanceSvc), authToken())
	res, err := client.ListDatasets(path)
	if err != nil {
		exitCode := 1
		switch err.(type) {
		case *maintenance.HTTPError:
			exitCode = 0
		default:
		}
		fmt.Fprintln(output, err.Error()+"\n")
		os.Exit(exitCode)
	} else if res != nil && len(*res) > 0 {
		printBeforeDelete(res, output, input, dryRun)
	} else {
		// no error and response is nil
		fmt.Fprintf(output, "Could not find any datasets to delete.\n\r")
	}

}

func printBeforeDelete(datasetsAndFolders *maintenance.ListDatasetResponse, output io.Writer, input io.Reader, dryRun bool) {
	for _, item := range *datasetsAndFolders {
		if item.IsDataset() {
			deleteWithPrompt(item.Path, output, input, dryRun)
		} else {
			deleteRecursively(item.Path, output, input, dryRun)
		}
	}
}

func deleteWithPrompt(path string, output io.Writer, input io.Reader, dryRun bool) {
	fmt.Fprint(output, "Delete dataset ", path, "? ")
	reader := bufio.NewReader(input)
	char, _, err := reader.ReadRune()

	if err != nil {
		fmt.Fprint(output, err)
	}

	switch char {
	case 'y', 'Y':
		doDelete(path, output, dryRun)
		break
	default:
		fmt.Fprintln(output, "... skipped")
	}
}

// Output:
// > dapla rm /foo/bar
//  Dataset /foo/bar (42 versions) successfully deleted
//
// > dapla rm --debug /foo/bar
//  /foo/bar [versiontimestamp]
//  	gs://bucket/random/prefix/file1 123
//  	gs://bucket/random/prefix/file2 123
//  	gs://bucket/random/prefix/file3 123
//  	gs://bucket/random/prefix/file4 123
//  /foo/bar [versiontimestamp2]
//  	gs://bucket2/some/other/path/prefix/file1 123
//  	gs://bucket2/some/other/path/prefix/file2 123
//  	gs://bucket2/some/other/path/prefix/file2 123
//  	gs://bucket2/some/other/path/prefix/file2 123
//  Number of deleted files: 1234
//  Total size: 123456789 KiB
//  Dataset /foo/bar (42 versions) successfully deleted
//
// > dapla rm --dry-run /foo/bar
//  Dataset /foo/bar (42 versions) successfully deleted
// The dry-run flag was set. NO FILES WERE DELETED.
func printDeleteResponse(deleteResponse *maintenance.DeleteDatasetResponse, output io.Writer, dryRun bool) {
	writer := bufio.NewWriter(output)
	defer writer.Flush()

	if viper.GetBool(CFGDebug) {
		for _, datasetVersion := range deleteResponse.DatasetVersion {
			fmt.Fprintf(writer, "Version: %s\n", datasetVersion.Timestamp)
			for _, deletedFile := range datasetVersion.DeletedFiles {
				fmt.Fprintf(writer, "\t%s\n", deletedFile.URI)
			}
		}
		fmt.Fprintf(writer, "\nnumber of deleted files: %d\n", deleteResponse.GetNumberOfFiles())
		fmt.Fprintf(writer, "total size of deleted files: %d\n", deleteResponse.TotalSize)
	}
	fmt.Fprintf(writer, "Dataset %s (%d %s) successfully deleted\n\r",
		deleteResponse.DatasetPath,
		len(deleteResponse.DatasetVersion),
		pluralize("version", len(deleteResponse.DatasetVersion)))

	if dryRun {
		fmt.Fprintf(writer, "The dry-run flag was set. NO FILES WERE DELETED.\n\r")
	}
}

func pluralize(text string, n int) string {
	if n > 1 {
		return text + "s"
	}
	return text
}
