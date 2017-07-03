package kong

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// TargetsService handles communication with Kong's '/targets' resource.
type TargetsService struct {
	*service
}

// Targets represents the object returned from Kong when querying for
// multiple target objects.
//
// In cases where the number of objects returned exceeds the maximum,
// Next holds the URI for the next set of results.
// i.e. "http://localhost:8001/targets/?size=2&offset=4d924084-1adb-40a5-c042-63b19db421d1"
type Targets struct {
	Data   []*Target `json:"data,omitempty"`
	Total  int       `json:"total,omitempty"`
	Next   string    `json:"next,omitempty"`
	Offset string    `json:"offset,omitempty"`
}

// TargetCount represents the number of targets on an upstream
type TargetCount struct {
	Total int `json:"total"`
}

// Target represents a single Kong target object.
type Target struct {
	Target     string `json:"target"`
	ID         string `json:"id,omitempty"`
	CreatedAt  int64  `json:"created_at,omitempty"`
	Weight     int    `json:"weight,omitempty"`
	UpstreamID string `json:"upstream_id,omitempty"`
}

// GetAllActive lists all the active targets attached to the specifed upstream.
//
// Equivalent to GET/upstreams/{name or id}/targets/active
func (s *TargetsService) GetAllActive(upstream string) (*Targets, *http.Response, error) {
	u := fmt.Sprintf("upstreams/%v/targets/active", upstream)

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	uResp := new(Targets)
	targetCount := new(TargetCount)
	// The targets/active endpoint has a bug where no results is manifest by data: {} instead of data: []
	// This is a workaround for this issue
	resp, err := s.client.Do(req, nil)
	if err != nil {
		return nil, resp, err
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, resp, err
	}

	err = json.NewDecoder(bytes.NewReader(respBody)).Decode(targetCount)
	if err == io.EOF {
		err = nil
	}
	if err != nil {
		return nil, resp, err
	}

	if targetCount.Total != 0 {
		// It is safe to cast to full Targets list
		err = json.NewDecoder(bytes.NewReader(respBody)).Decode(uResp)
	} else {
		uResp.Data = []*Target{}
		uResp.Total = 0
	}
	if err == io.EOF {
		err = nil
	}

	if err != nil {
		return nil, resp, err
	}

	return uResp, resp, err
}

// Delete deletes a single Kong target object, by target (host/port combination)
//
// Equivalent to DELETE /upstreams/{name or id}/targets/{target}
func (s *TargetsService) Delete(upstream string, target string) (*http.Response, error) {
	u := fmt.Sprintf("upstreams/%v/targets/%v", upstream, target)

	req, err := s.client.NewRequest(http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, err
}

// Post creates a new Kong target object.
//
// Equivalent to POST /upstreams/{name or id}/targets
func (s *TargetsService) Post(upstream string, target *Target) (*http.Response, error) {
	req, err := s.client.NewRequest(http.MethodPost, fmt.Sprintf("upstreams/%v/targets", upstream), target)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)

	return resp, err
}
