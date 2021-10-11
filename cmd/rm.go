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
				deleteRecursively(path, rmDryRun)
			} else {
				doDelete(path, rmDryRun)
			}
		},
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {

			return doAutoComplete(toComplete)

		},
	}
}

func doDelete(path string, dryRun bool) {
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
		fmt.Println(err.Error() + "\n")
		os.Exit(exitCode)
	} else {
		printDeleteResponse(res, os.Stdout, dryRun)
	}
}

func deleteRecursively(path string, dryRun bool) {
	var client = maintenance.NewClient(apiURLOf(APINameDataMaintenanceSvc), authToken())
	res, err := client.ListDatasets(path)
	if err != nil {
		exitCode := 1
		switch err.(type) {
		case *maintenance.HTTPError:
			exitCode = 0
		default:
		}
		fmt.Println(err.Error() + "\n")
		os.Exit(exitCode)
	} else if res != nil {
		printBeforeDelete(res, dryRun)
	} else {
		// no error and response is nil
		writer := bufio.NewWriter(os.Stdout)
		fmt.Fprintf(writer, "Could not find any datasets to delete.\n\r")
	}

}

func printBeforeDelete(datasetsAndFolders *maintenance.ListDatasetResponse, dryRun bool) {
	for _, item := range *datasetsAndFolders {
		if item.IsDataset() {
			deleteWithPrompt(item.Path, dryRun)
		} else {
			deleteRecursively(item.Path, dryRun)
		}
	}
}

func deleteWithPrompt(path string, dryRun bool) {
	fmt.Print("Delete dataset ", path, "? ")
	reader := bufio.NewReader(os.Stdin)
	char, _, err := reader.ReadRune()

	if err != nil {
		fmt.Println(err)
	}

	switch char {
	case 'Y':
	case 'y':
		doDelete(path, dryRun)
		break
	default:
		fmt.Println("... skipped")
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
