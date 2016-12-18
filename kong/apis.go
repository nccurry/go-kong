package kong

import (
	"errors"
	"fmt"
	"net/http"
)

// ApisService handles communication with the kong '/apis' resource.
type ApisService service

// Apis represents the object returned from Kong when querying for multiple api objects.
//
// In cases where the number of objects returned exceeds the maximum, Next holds the
// URI for the next set of results.
// i.e. "http://localhost:8001/apis/?size=2&offset=4d924084-1adb-40a5-c042-63b19db421d1"
type Apis struct {
	Data  []Api  `json:"data,omitempty"`
	Total int    `json:"total,omitempty"`
	Next  string `json:"next,omitempty"`
}

// Api represents a Kong api object.
type Api struct {
	UpstreamURL      string `json:"upstream_url,omitempty"`
	StripRequestPath bool   `json:"strip_request_path,omitempty"`
	RequestPath      string `json:"request_path,omitempty"`
	ID               string `json:"id,omitempty"`
	CreatedAt        int64  `json:"created_at,omitempty"`
	PreserveHost     bool   `json:"preserve_host,omitempty"`
	Name             string `json:"name,omitempty"`
}

// ApisService.Get queries for a single kong api object, by name or id.
// Equivalent to GET /apis/{name or id}
//
// In addition to the *Api object, the *http.Response from kong
// is returned if the caller wishes to do further inspection.
func (s *ApisService) Get(api string) (*Api, *http.Response, error) {
	u := fmt.Sprintf("apis/%v", api)

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

// ApisService.Patch updates an existing kong api object.
// At least one of api.Name or api.ID must be specified in the passed *Api parameter.
// Equivalent to PATCH /apis/{name or id}
//
// The *http.Response from kong is returned if the caller wishes to do further inspection.
func (s *ApisService) Patch(api *Api) (*http.Response, error) {
	var u string
	if api.Name != "" {
		u = fmt.Sprintf("apis/%v", api.Name)
	} else if api.ID != "" {
		u = fmt.Sprintf("apis/%v", api.ID)
	} else {
		return nil, errors.New("At least one of api.Name or api.ID must be specified")
	}

	req, err := s.client.NewRequest("PATCH", u, api)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)

	return resp, err

}

// ApisService.Delete deletes a single kong api object, by name or id.
// Equivalent to DELETE /apis/{name or id}
//
// The *http.Response from kong is returned if the caller wishes to do further inspection.
func (s *ApisService) Delete(api string) (*http.Response, error) {
	u := fmt.Sprintf("apis/%v", api)

	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, err
}

// ApisService.Post cretes a new kong api object.
// Equivalent to POST /apis
//
// The *http.Response from kong is returned if the caller wishes to do further inspection.
func (s *ApisService) Post(api *Api) (*http.Response, error) {
	req, err := s.client.NewRequest("POST", "apis", api)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)

	return resp, err
}

// ApisGetAllOptions specifies optional filter parameters to the ApisService.GetAll method.
//
// Additional information about filtering options can be found in the Kong documentation at
// https://getkong.org/docs/0.9.x/admin-api/#list-apis
type ApisGetAllOptions struct {
	ID          string `url:"id,omitempty"`           // A filter on the list based on the apis id field.
	Name        string `url:"name,omitempty"`         // A filter on the list based on the apis name field.
	RequestHost string `url:"request_host,omitempty"` // A filter on the list based on the apis request_host field.
	RequestPath string `url:"request_path,omitempty"` // A filter on the list based on the apis request_path field.
	UpstreamURL string `url:"upstream_url,omitempty"` // A filter on the list based on the apis upstream_url field.
	Size        int    `url:"size,omitempty"`         // A limit on the number of objects to be returned.
	Offset      string `url:"offset,omitempty"`       // A cursor used for pagination. offset is an object identifier that defines a place in the list.
}

// ApisService.GetAll queries for all Kong api objects.
// This query can be filtered by supplying the ApisGetAllOptions struct.
// Equivalent to GET /apis?uri=params&from=opt
//
// In addition to the *Apis object, the *http.Response from kong
// is returned if the caller wishes to do further inspection.
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
