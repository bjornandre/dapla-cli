package rest

import (
	"github.com/h2non/gock"
	"net/http"
	"testing"
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
    "name": "foo/file1"
},{
    "createdBy": "Kari Nordmann",
    "createdDate": "3000-01-01T00:00:00.123456Z",
    "name": "foo/file2"
}]
`)

	gock.New("http://server.com").
		Reply(http.StatusForbidden)

	var client = NewClient("http://server.com", "a secret secret")

	datasets, err := client.ListDatasets("foo")
	if err != nil {
		t.Errorf("Error %v", err)
	} else if len(*datasets) != 2 {
		t.Errorf("Invalid response")
	}
}
