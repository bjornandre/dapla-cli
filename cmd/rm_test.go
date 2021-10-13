package cmd

import (
	"bytes"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/andreyvit/diff"
	"github.com/spf13/viper"
	"github.com/statisticsnorway/dapla-cli/maintenance"
	"gopkg.in/h2non/gock.v1"
)

func TestExecuteRM(t *testing.T) {

	tests := []struct {
		response                     maintenance.DeleteDatasetResponse
		expectedOutput               string
		expectedOutputDebug          string
		expectedOutputDryRun         string
		expectedOutputDebugAndDryRun string
	}{
		{response: maintenance.DeleteDatasetResponse{
			DatasetPath: "/foo/bar",
			TotalSize:   15,
			DatasetVersion: []maintenance.DatasetVersion{
				{
					Timestamp: time.Date(2000, 1, 1, 0, 0, 0, 123456000, time.UTC),
					DeletedFiles: []maintenance.DatasetFile{
						{URI: "gs://bucket/prefix/foo/bar/v1/file1", Size: 1},
						{URI: "gs://bucket/prefix/foo/bar/v1/file2", Size: 2},
					},
				},
				{
					Timestamp: time.Date(3000, 1, 1, 0, 0, 0, 123456000, time.UTC),
					DeletedFiles: []maintenance.DatasetFile{
						{URI: "gs://bucket/prefix/foo/bar/v2/file1", Size: 4},
						{URI: "gs://bucket/prefix/foo/bar/v2/file2", Size: 8},
					},
				},
			},
		},

			expectedOutput: "Dataset /foo/bar (2 versions) successfully deleted",
			expectedOutputDebug: "Version: 2000-01-01 00:00:00.123456 +0000 UTC\n" +
				"\tgs://bucket/prefix/foo/bar/v1/file1\n" +
				"\tgs://bucket/prefix/foo/bar/v1/file2\n" +
				"Version: 3000-01-01 00:00:00.123456 +0000 UTC\n" +
				"\tgs://bucket/prefix/foo/bar/v2/file1\n" +
				"\tgs://bucket/prefix/foo/bar/v2/file2\n\n" +
				"number of deleted files: 4\n" +
				"total size of deleted files: 15\n" +
				"Dataset /foo/bar (2 versions) successfully deleted",
			expectedOutputDryRun: "Dataset /foo/bar (2 versions) successfully deleted\n\r" +
				"The dry-run flag was set. NO FILES WERE DELETED.",
			expectedOutputDebugAndDryRun: "Version: 2000-01-01 00:00:00.123456 +0000 UTC\n" +
				"\tgs://bucket/prefix/foo/bar/v1/file1\n" +
				"\tgs://bucket/prefix/foo/bar/v1/file2\n" +
				"Version: 3000-01-01 00:00:00.123456 +0000 UTC\n" +
				"\tgs://bucket/prefix/foo/bar/v2/file1\n" +
				"\tgs://bucket/prefix/foo/bar/v2/file2\n\n" +
				"number of deleted files: 4\n" +
				"total size of deleted files: 15\n" +
				"Dataset /foo/bar (2 versions) successfully deleted\n\r" +
				"The dry-run flag was set. NO FILES WERE DELETED.",
		},
	}

	for _, values := range tests {
		var output bytes.Buffer

		// Test rm without flags
		printDeleteResponse(&values.response, &output, false)
		if actual, expected := strings.TrimSpace(output.String()),
			strings.TrimSpace(values.expectedOutput); actual != expected {
			fmt.Println("***** <rm> WITHOUT FLAGS *****")
			t.Errorf("Result not as expected:\n%v", diff.LineDiff(expected, actual))
		}
		output.Reset()

		// Test rm with debug flag
		viper.Set(CFGDebug, true)
		printDeleteResponse(&values.response, &output, false)

		if actual, expected := strings.TrimSpace(output.String()),
			strings.TrimSpace(values.expectedOutputDebug); actual != expected {
			fmt.Println("***** <rm> WITH DEBUG FLAG *****")
			t.Errorf("Result not as expected:\n%v", diff.LineDiff(expected, actual))
		}
		output.Reset()

		// Test rm with dry-run flag
		viper.Set(CFGDebug, false)
		printDeleteResponse(&values.response, &output, true)

		if actual, expected := strings.TrimSpace(output.String()),
			strings.TrimSpace(values.expectedOutputDryRun); actual != expected {
			fmt.Println("***** <rm> WITH DRY-RUN FLAG *****")
			t.Errorf("Result not as expected:\n%v", diff.LineDiff(expected, actual))
		}
		output.Reset()

		// Test rm with both debug and dry-run flags
		viper.Set(CFGDebug, true)
		printDeleteResponse(&values.response, &output, true)

		if actual, expected := strings.TrimSpace(output.String()),
			strings.TrimSpace(values.expectedOutputDebugAndDryRun); actual != expected {
			fmt.Println("***** <rm> WITH DEBUG AND DRY-RUN FLAG *****")
			t.Errorf("Result not as expected:\n%v", diff.LineDiff(expected, actual))
		}
		output.Reset()

	}
}

// Note that the "spinner" makes IntelliJ show 'No tests were run' (https://youtrack.jetbrains.com/issue/GO-7215)
func TestDelete(t *testing.T) {
	readTestConfig()
	defer gock.Off()

	gock.New("http://data-maintenance-mock/(.*)").
		PathParam("delete", "foo").
		MatchParam("dry-run", "false").
		MatchHeader("Authorization", "Bearer eyJh...TqV2Q").
		Reply(http.StatusOK).
		BodyString(`{
			"datasetPath": "foo",
			"deletedVersions": [{
				"timestamp": "2012-04-23T18:25:43.511Z"
			}],
			"totalSize": 10
		}`)

	output, err := executeCommand(newRmCommand(), "foo")

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	assert.Equal(t, "Dataset foo (1 version) successfully deleted\n\r", output)
	assert.True(t, gock.IsDone())
}

func TestDeleteRecursivelyNoData(t *testing.T) {
	readTestConfig()
	defer gock.Off()

	gock.New("http://data-maintenance-mock/(.*)").
		PathParam("list", "foo").
		MatchHeader("Authorization", "Bearer eyJh...TqV2Q").
		Reply(http.StatusOK).
		BodyString("[]")

	command := newRmCommand()
	command.Flags().BoolVarP(&rmRecursive, "recursive", "", false, "delete recursively")
	output, err := executeCommand(command, "--recursive", "foo")

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	t.Logf(output)
	assert.Equal(t, "Could not find any datasets to delete.\n\r", output)
	assert.True(t, gock.IsDone())
	gock.Clean()
}

func TestDeleteRecursively(t *testing.T) {
	readTestConfig()
	defer gock.Off()

	gock.New("http://data-maintenance-mock/(.*)").
		Path("list/foo").
		MatchHeader("Authorization", "Bearer eyJh...TqV2Q").
		Reply(http.StatusOK).
		BodyString(`[{
			"path": "foo/bar",
			"createdDate": "2012-04-23T18:25:43.511Z",
			"depth": 1
		}]`)

	gock.New("http://data-maintenance-mock/(.*)").
		Path("list/foo/bar").
		MatchHeader("Authorization", "Bearer eyJh...TqV2Q").
		Reply(http.StatusOK).
		BodyString(`[{
			"path": "foo/bar/baz",
			"createdDate": "2012-04-23T18:25:43.511Z",
			"depth": 0
		}]`)

	gock.New("http://data-maintenance-mock/(.*)").
		Path("delete/foo/bar/baz").
		MatchParam("dry-run", "false").
		MatchHeader("Authorization", "Bearer eyJh...TqV2Q").
		Reply(http.StatusOK).
		BodyString(`{
			"datasetPath": "foo/bar/baz",
			"deletedVersions": [{
				"timestamp": "2012-04-23T18:25:43.511Z"
			}],
			"totalSize": 10
		}`)

	command := newRmCommand()
	command.Flags().BoolVarP(&rmRecursive, "recursive", "", false, "delete recursively")

	var stdin bytes.Buffer
	// Simulate that the user enters 'y' from stdin
	stdin.Write([]byte("y"))
	command.SetIn(&stdin)
	output, err := executeCommand(command, "--recursive", "foo")

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	t.Logf(output)
	assert.Equal(t, "Delete dataset foo/bar/baz? Dataset foo/bar/baz (1 version) successfully deleted\n\r", output)
	assert.True(t, gock.IsDone())
	gock.Clean()
}

func TestDeleteRecursivelyAnswerNo(t *testing.T) {
	readTestConfig()
	defer gock.Off()

	gock.New("http://data-maintenance-mock/(.*)").
		Path("list/foo").
		MatchHeader("Authorization", "Bearer eyJh...TqV2Q").
		Reply(http.StatusOK).
		BodyString(`[{
			"path": "foo/bar",
			"createdDate": "2012-04-23T18:25:43.511Z",
			"depth": 0
		}]`)

	command := newRmCommand()
	command.Flags().BoolVarP(&rmRecursive, "recursive", "", false, "delete recursively")

	var stdin bytes.Buffer
	// Simulate that the user enters 'n' from stdin
	stdin.Write([]byte("n"))
	command.SetIn(&stdin)
	output, err := executeCommand(command, "--recursive", "foo")

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	t.Logf(output)
	assert.Equal(t, "Delete dataset foo/bar? ... skipped\n", output)
	assert.True(t, gock.IsDone())
	gock.Clean()
}

func executeCommand(root *cobra.Command, args ...string) (output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetArgs(args)

	_, err = root.ExecuteC()

	return buf.String(), err
}

func readTestConfig() {
	viper.SetConfigFile(".dapla-cli-tests.yml")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}
}
