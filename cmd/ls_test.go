package cmd

import (
	"bytes"
	"fmt"
	"github.com/statisticsnorway/dapla-cli/rest"
	"testing"
)

func TestExecute(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Execute(); (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test(t *testing.T) {

	tests := []struct {
		response       rest.ListDatasetResponse
		expectedOutput string
	}{
		{rest.ListDatasetResponse{
			rest.ListDatasetElement{Path: "/foo/bar"},
			rest.ListDatasetElement{Path: "/foo/baz"},
		},
			"/foo/bar\n\r/foo/baz",
		},
		{rest.ListDatasetResponse{
			rest.ListDatasetElement{Path: "/foo2/bar"},
			rest.ListDatasetElement{Path: "/foo2/baz"},
		},
			"/foo2/bar\n\r/foo2/baz",
		},
	}

	for _, values := range tests {
		var output bytes.Buffer
		printNewLine(&values.response, &output)

		if output.String() == values.expectedOutput {
			fmt.Println(output.String())
			t.Errorf("Invalid output, expected %v, got %v", values.expectedOutput, output.String())
		}
	}
}
