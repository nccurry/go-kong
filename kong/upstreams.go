package kong

import (
	"errors"
	"fmt"
	"net/http"
)

// UpstreamsService handles communication with Kong's '/upstreams' resource.
type UpstreamsService struct {
	*service
}

// Upstreams represents the object returned from Kong when querying for
// multiple upstream objects.
//
// In cases where the number of objects returned exceeds the maximum,
// Next holds the URI for the next set of results.
// i.e. "http://localhost:8001/upstreams/?size=2&offset=4d924084-1adb-40a5-c042-63b19db421d1"
type Upstreams struct {
	Data   []*Upstream `json:"data,omitempty"`
	Total  int         `json:"total,omitempty"`
	Next   string      `json:"next,omitempty"`
	Offset string      `json:"offset,omitempty"`
}

// Upstream represents a single Kong upstream object.
type Upstream struct {
	Name      string `json:"name"`
	ID        string `json:"id,omitempty"`
	CreatedAt int64  `json:"created_at,omitempty"`
	Slots     int    `json:"slots,omitempty"`
	Orderlist []int  `json:"orderliste,omitempty"`
}

// Get queries for a single Kong upstream object, by name or id.
//
// Equivalent to GET /upstreams/{name or id}
func (s *UpstreamsService) Get(upstream string) (*Upstream, *http.Response, error) {
	u := fmt.Sprintf("upstreams/%v", upstream)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	uResp := new(Upstream)
	resp, err := s.client.Do(req, uResp)
	if err != nil {
		return nil, resp, err
	}

	return uResp, resp, err
}

// Patch updates an existing Kong upstream object.
// At least one of upstream.Name or upstream.ID must be specified in
// the passed *Upstream parameter.
//
// Equivalent to PATCH /upstreams/{name or id}
func (s *UpstreamsService) Patch(upstream *Upstream) (*http.Response, error) {
	var u string
	if upstream.Name != "" {
		u = fmt.Sprintf("upstreams/%v", upstream.Name)
	} else if upstream.ID != "" {
		u = fmt.Sprintf("upstreams/%v", upstream.ID)
	} else {
		return nil, errors.New("At least one of upstream.Name or upstream.ID must be specified")
	}

	req, err := s.client.NewRequest("PATCH", u, upstream)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)

	return resp, err

}

// Delete deletes a single Kong upstream object, by name or id.
//
// Equivalent to DELETE /upstreams/{name or id}
func (s *UpstreamsService) Delete(upstream string) (*http.Response, error) {
	u := fmt.Sprintf("upstreams/%v", upstream)

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

// Post creates a new Kong upstream object.
//
// Equivalent to POST /upstreams
func (s *UpstreamsService) Post(upstream *Upstream) (*http.Response, error) {
	req, err := s.client.NewRequest("POST", "upstreams", upstream)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)

	return resp, err
}
