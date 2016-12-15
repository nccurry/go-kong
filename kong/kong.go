package kong

import (
	"net/http"
	"net/url"
	"io"
	"bytes"
	"encoding/json"
)

const (
	contentType = "Content-Type"
	applicationJson = "application/json"
)

// A Client manages communication with the Kong API
type Client struct {
	client *http.Client // HTTP client used to communicate with the API

	// Base URL for API requests.
	// BaseURL should always be specified with a trailing slash
	BaseURL *url.URL

        // Reuse a single struct instead of allocating one for each service on the heap
	common service

	// Services used for talking to different parts of the Kong API
	Apis *ApisService

}

type service struct {
	client *Client
}

func NewClient(httpClient *http.Client, baseURLStr string) (*Client, err) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	baseURL, err := url.Parse(baseURLStr)
	if err != nil {
		return nil, err
	}

	c := &Client{client: httpClient, BaseURL: baseURL}
	c.common.client = c
	c.Apis = (*ApisService)(&c.common)

	return c, nil
}

func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	req.Header.Set("Accept", "application/json")

	return req, nil
}