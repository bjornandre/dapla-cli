package rest

import (
	"github.com/h2non/gock"
	"github.com/steinfletcher/apitest"
	"net/http"
	"testing"
)

var listDataset = apitest.NewMock().
	Get("/api/v1/list/foo/bar").
	RespondWith().
	Body(`
[{
	"createdBy": "Arild Johan Takvam-Borge",
	"createdDate": "2020-11-12T20:52:00.528414Z",
    "name": "skatt.tmp"
},{
    "createdBy": "Hadrien Kohl",
    "createdDate": "2020-11-12T20:35:20.528414Z",
    "name": "teststuff"
}]`).
	Status(http.StatusOK).
	End()

func TestClient_DeleteDatasets(t *testing.T) {
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
