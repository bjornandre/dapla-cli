package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	BaseURL    string
	Client     *http.Client
	authBearer string
}

type DatasetResponse []DatasetElement

type DatasetElement struct {
	Name      string    `json:"name"`
	CreatedBy string    `json:"createdBy"`
	CreatedAt time.Time `json:"createdDate"`
}

func NewClient(baseURL string, authBearer string) *Client {
	return &Client{
		BaseURL:    baseURL,
		Client:     http.DefaultClient,
		authBearer: authBearer,
	}
}

func NewClientWithJupyter(baseURL string) *Client {
	// TODO: Get jupyter URL from env.
	// 	     Get token from env
	// 		 Fetch token.
	return &Client{
		BaseURL:    baseURL,
		Client:     http.DefaultClient,
		authBearer: "TODO",
	}
}

func (c Client) DeleteDatasets(path string) error {
	panic("TODO")
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
