package rest

import (
	"github.com/h2non/gock"
	"net/http"
	"testing"
	"time"
)

func TestClient_fetchJupyterToken(t *testing.T) {
	defer gock.Off()

	gock.New("http://server.com").
		Get("/foo/bar/token").
		MatchHeader("Authorization", "^token the api token$").
		Reply(http.StatusOK).
		BodyString(`{ "access_token": "the access token"}`)

	gock.New("http://server.com").
		Reply(http.StatusForbidden)

	token, err := fetchJupyterToken("http://server.com/foo/bar/token", "the api token")
	if err != nil {
		t.Fatal(err)
	}

	if token != "the access token" {
		t.Errorf("expected %s but got %s", "the access token", token)
	}
}

func TestClient_ListDatasets(t *testing.T) {
	defer gock.Off()

	gock.New("http://server.com").
		Get("/api/v1/list/foo").
		MatchHeader("Authorization", "^Bearer a secret secret$").
		Reply(http.StatusOK).BodyString(`
[{
	"createdBy": "Ola Nordmann",
	"createdDate": "2000-01-01T00:00:00.123456Z",
    "path": "foo/file1",
	"type": "BOUNDED",
	"valuation": "INTERNAL",
	"state": "INPUT"
},{
    "createdBy": "Kari Nordmann",
    "createdDate": "3000-01-01T00:00:00.123456Z",
    "path": "foo/file2",
	"type": "unBOUNDED",
	"valuation": "SENSITIVE",
	"state": "RAW"
}]
`)

	gock.New("http://server.com").
		Reply(http.StatusForbidden)

	var client = NewClient("http://server.com", "a secret secret")

	datasets, err := client.ListDatasets("foo")
	// TODO use these in the reply above
	expectedCreatedBy := "Ola Nordmann"
	expectedCreatedAt := "2000-01-01T00:00:00.123456Z"
	expectedPath := "foo/file1"
	expectedType := "BOUNDED"
	expectedValuation := "INTERNAL"
	expectedState := "INPUT"
	var v = *datasets
	if v[0].CreatedBy != expectedCreatedBy {
		t.Errorf("Expected %v, but got %v", expectedCreatedBy, v[0].CreatedBy)
	}

	//TODO how to compare these correctly?
	parsedExpectedCreatedAt, _ := time.Parse(expectedCreatedAt, expectedCreatedAt)
	if v[0].CreatedAt != parsedExpectedCreatedAt {
		t.Errorf("Expected %v, but got %v", expectedCreatedAt, v[0].CreatedAt)
	}
	if v[0].Path != expectedPath {
		t.Errorf("Expected %v, but got %v", expectedPath, v[0].Path)
	}
	if v[0].Type != expectedType {
		t.Errorf("Expected %v, but got %v", expectedType, v[0].Type)
	}
	if v[0].Valuation != expectedValuation {
		t.Errorf("Expected %v, but got %v", expectedValuation, v[0].Valuation)
	}
	if v[0].State != expectedState {
		t.Errorf("Expected %v, but got %v", expectedState, v[0].State)
	}
	if err != nil {
		t.Errorf("Error %v", err)
	} else if len(*datasets) != 2 {
		t.Errorf("Invalid response")
	}
}
