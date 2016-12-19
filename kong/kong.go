package kong

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/go-querystring/query"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
)

const (
	contentType     = "Content-Type"
	applicationJson = "application/json"
)

// Client manages communication with the Kong API.
// New client objects should be created using the NewClient function.
// The BaseURL field must be defined and pointed at an instance of the
// Kong Admin API.
//
// Kong resources can be access using the Service objects.
// client.Apis.Get("id") -> GET /apis/id
// client.Consumers.Patch(consumer) -> PATCH /consumers/id
type Client struct {
	client *http.Client // HTTP client used to communicate with the API

	// Base URL for API requests.
	// BaseURL should always be specified with a trailing slash
	BaseURL *url.URL

	// Reuse a single struct instead of allocating one for each service on the heap
	common service

	// Services used for talking to different parts of the Kong API
	Apis      *ApisService
	Consumers *ConsumersService
	Plugins   *PluginsService
}

// Each service representing a Kong resource type will be of this type
type service struct {
	client *Client
}

// addOptions is used to convert the opt struct passed to the
// various Service.GetAll(opt) methods to query string parameters
// on the end of the resource URI passed for s. It is assumed that
// the struct fields passed in for opt are tagged with `url:""` tags.
//
// i.e &ApisGetAllOptions{RequestPath: "service"} -> /apis?request_path=service
func addOptions(s string, opt interface{}) (string, error) {
	v := reflect.ValueOf(opt)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	qs, err := query.Values(opt)
	if err != nil {
		return s, err
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}

// NewClient creates a new kong.Client object.
// This should be the primary way a kong.Client object is constructed.
//
// If an httpClient object is specified it will be used instead of the
// default http.DefaultClient.
//
// baseURLStr should point to an instance a Kong Admin API and must
// contain the trailing slash. i.e. http://kong:8001/
func NewClient(httpClient *http.Client, baseURLStr string) (*Client, error) {
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
	c.Consumers = (*ConsumersService)(&c.common)
	c.Plugins = (*PluginsService)(&c.common)

	return c, nil
}

// NewRequest is used to construct a new *http.Request object
// Generally speaking the returned *http.Request object will be
// used in a subsequent Client.Do to execute the actual REST call
// against Kong.
//
// If body is provided, it will be JSON encoded and used as the request
// body.
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

// Do executes the actual REST call against Kong. The API response is
// JSON decoded and stored in the value pointed to by v, or returned as an
// error if an API error has occurred.  If v implements the io.Writer
// interface, the raw response body will be written to v, without attempting to
// first decode it.
//
// The *http.Response object returned by Do should eventually get
// passed back to the caller. If Kong returns a status code outside
// of the 200 range, the caller can inspect the *http.Response to
// get more information. Additionally the err returned in this case
// will be of type ErrorResponse.
func (c *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	err = CheckResponse(resp)
	if err != nil {
		// even though there was an error, we still return the response
		// in case the caller wants to inspect it further
		return resp, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
			if err == io.EOF {
				err = nil // ignore EOF errors caused by empty response body
			}
		}
	}

	return resp, err
}

// ErrorResponse is returned from Client.Do if Kong returns a status
// code outside the 200 range.
type ErrorResponse struct {
	Response    *http.Response // HTTP response that caused this error
	KongMessage string         `json:"message,omitempty"`
	KongError   string         `json:"error,omitempty"`
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %v %v",
		r.Response.Request.Method, r.Response.Request.URL,
		r.Response.StatusCode, r.KongMessage, r.KongError)
}

// ConflictError occurs when trying to create a resource that already exists.
// CheckResponse will return this type of error when Kong returns a 409 status code.
type ConflictError ErrorResponse

func (r *ConflictError) Error() string {
	return (*ErrorResponse)(r).Error()
}

// NotFoundError occurs when trying to access a resource that does not exist.
// CheckResponse will return this type of error when Kong returns a 404 status code.
type NotFoundError ErrorResponse

func (r *NotFoundError) Error() string {
	return (*ErrorResponse)(r).Error()
}

// CheckResponse looks at the response from a Kong API call
// and determines what type of error needs to be returned, if any.
func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}
	errorResponse := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && data != nil {
		json.Unmarshal(data, errorResponse)
	}

	switch r.StatusCode {
	case 404:
		return (*NotFoundError)(errorResponse)
	case 409:
		return (*ConflictError)(errorResponse)
	default:
		return errorResponse
	}
}
