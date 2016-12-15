package kong

import (
	"net/http"
	"net/url"
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