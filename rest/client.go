package rest

import (
	"encoding/json"
	"fmt"
	errors2 "github.com/pkg/errors"
	"net/http"
	"net/url"
	"os"
	"time"
)

type Client struct {
	BaseURL    string
	Client     *http.Client
	authBearer string
}

type DatasetResponse []DatasetElement

const jupyterHUBTokenURL = "JUPYTERHUB_HANDLER_CUSTOM_AUTH_URL"
const jupyterAPIToken = "JUPYTERHUB_API_TOKEN"

type DatasetElement struct {
	Path      string    `json:"path"`
	CreatedBy string    `json:"createdBy"`
	CreatedAt time.Time `json:"createdDate"`
	Type      string    `json:"type"`
	Valuation string    `json:"valuation"`
	State     string    `json:"state"`
}

func NewClient(baseURL string, authBearer string) *Client {
	return &Client{
		BaseURL:    baseURL,
		Client:     http.DefaultClient,
		authBearer: authBearer,
	}
}

func NewClientWithJupyter(baseURL string) (*Client, error) {

	apiURL := os.Getenv(jupyterHUBTokenURL)
	apiToken := os.Getenv(jupyterAPIToken)
	if apiToken == "" || apiURL == "" {
		return nil, errors2.Errorf("missing environment %s or %s", jupyterHUBTokenURL, jupyterAPIToken)
	}

	token, err := fetchJupyterToken(apiURL, apiToken)
	if err != nil {
		return nil, err
	}
	return &Client{
		BaseURL:    baseURL,
		Client:     http.DefaultClient,
		authBearer: token,
	}, nil
}

// Fetch the JTW token from jupyter environment
func fetchJupyterToken(apiURL, apiToken string) (string, error) {
	parsedURL, err := url.Parse(apiURL)
	if err != nil {
		return "", err
	}

	client := http.DefaultClient
	req, err := http.NewRequest(http.MethodGet, parsedURL.String(), nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", fmt.Sprintf("token %s", apiToken))
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		return "", fmt.Errorf("unknown error, status code: %d", res.StatusCode)
	}

	var data map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return "", nil
	}

	return data["access_token"].(string), nil
}

func (c Client) DeleteDatasets(path string) error {
	panic(fmt.Sprintf("TODO %s", path))
}

func (c *Client) ListDatasets(path string) (*DatasetResponse, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/list/%s", c.BaseURL, path), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.authBearer))
	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	res, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("unknown error, status code: %d", res.StatusCode)
	}

	resp := DatasetResponse{}
	err = json.NewDecoder(res.Body).Decode(&resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
