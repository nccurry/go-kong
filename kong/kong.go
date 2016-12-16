package kong

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"github.com/google/go-querystring/query"
	"io/ioutil"
	"log"
)

const (
	contentType     = "Content-Type"
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
	Consumers *ConsumersService
	Plugins *PluginsService
}

type service struct {
	client *Client
}

// addOptions adds the parameters in opt as URL query parameters to s.  opt
// must be a struct whose fields may contain "url" tags.
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

	b, _ := ioutil.ReadAll(buf)
	log.Println(string(b))

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

// Do sends an API request and returns the API response.  The API response is
// JSON decoded and stored in the value pointed to by v, or returned as an
// error if an API error has occurred.  If v implements the io.Writer
// interface, the raw response body will be written to v, without attempting to
// first decode it.
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

func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}
	errorResponse := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && data != nil {
		json.Unmarshal(data, errorResponse)
	}

	return errorResponse // TODO: Return other kinds of errors
}
