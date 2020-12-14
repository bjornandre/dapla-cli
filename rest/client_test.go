package rest

import (
	"github.com/google/go-cmp/cmp"
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
	"type": "UNBOUNDED",
	"valuation": "SENSITIVE",
	"state": "RAW"
}]
`)

	gock.New("http://server.com").
		Reply(http.StatusForbidden)

	var client = NewClient("http://server.com", "a secret secret")

	var expected = DatasetElement{
		CreatedAt: time.Date(2000, 1, 1, 0, 0, 0, 123456000, time.UTC),
		CreatedBy: "Ola Nordmann",
		Path:      "foo/file1",
		Type:      "BOUNDED",
		Valuation: "INTERNAL",
		State:     "INPUT",
	}

	datasets, err := client.ListDatasets("foo")
	if err != nil {
		t.Errorf("Got error %v", err)
	}

	if len(*datasets) != 2 {
		t.Errorf("Invalid response")
	}
	var element = (*datasets)[0]
	if !cmp.Equal(expected, element) {
		t.Errorf("Expected %v, but got %v", expected, element)
	}
}
