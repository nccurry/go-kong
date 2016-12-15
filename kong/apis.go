package kong

import (
	"fmt"
	"net/http"
)

// ApisService handles communication with the apis resource of Kong
type ApisService service

// Apis represents a list of Kong Apis
type Apis struct {
	Data  []Api  `json:"data,omitempty"`
	Total int    `json:"total,omitempty"`
	Next  string `json:"next,omitempty"`
}

type Api struct {
	UpstreamURL      string `json:"upstream_url,omitempty"`
	StripRequestPath bool   `json:"strip_request_path,omitempty"`
	RequestPath      string `json:"request_path,omitempty"`
	ID               string `json:"id,omitempty"`
	CreatedAt        int64  `json:"created_at,omitempty"`
	PreserveHost     bool   `json:"preserve_host,omitempty"`
	Name             string `json:"name,omitempty"`
}

func (s *ApisService) Get(apiName string) (*Api, *http.Response, error) {
	u := fmt.Sprintf("apis/%v", apiName)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	uResp := new(Api)
	resp, err := s.client.Do(req, uResp)
	if err != nil {
		return nil, resp, err
	}

	return uResp, resp, err
}

func (s *ApisService) Patch(api *Api) (*http.Response, error) {
	u := fmt.Sprintf("apis/%v", api.Name)

	req, err := s.client.NewRequest("PATCH", u, api)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)

	return resp, err

}

func (s *ApisService) Post(api *Api) (*http.Response, error) {
	req, err := s.client.NewRequest("POST", "apis", api)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)

	return resp, err
}

// ApisGetAllOptions specifies optional filter parameters to the ApisService.GetAll method
type ApisGetAllOptions struct {
	ID          string `url:"id,omitempty"`           // A filter on the list based on the apis id field.
	Name        string `url:"name,omitempty"`         // A filter on the list based on the apis name field.
	RequestHost string `url:"request_host,omitempty"` // A filter on the list based on the apis request_host field.
	RequestPath string `url:"request_path,omitempty"` // A filter on the list based on the apis request_path field.
	UpstreamURL string `url:"upstream_url,omitempty"` // A filter on the list based on the apis upstream_url field.
	Size        int    `url:"size,omitempty"`         // A limit on the number of objects to be returned.
	Offset      string `url:"offset,omitempty"`       // A cursor used for pagination. offset is an object identifier that defines a place in the list.
}

// GetAll lists all apis
func (s *ApisService) GetAll(opt *ApisGetAllOptions) (*Apis, *http.Response, error) {
	u, err := addOptions("apis", opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	apis := new(Apis)
	resp, err := s.client.Do(req, apis)
	if err != nil {
		return nil, resp, err
	}

	return apis, resp, err
}
