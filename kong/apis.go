package kong

import (
	"net/http"
	"fmt"
)


// ApisService handles communication with the apis resource of Kong
type ApisService service

// Apis represents a list of Kong Apis
type Apis struct {
	Data []Api `json:"data,omitempty"`
	Total int `json:"total,omitempty"`
	Next string `json:"next,omitempty"`
}

type Api struct {
	UpstreamURL      string `json:"upstream_url"`
	StripRequestPath bool   `json:"strip_request_path"`
	RequestPath      string `json:"request_path"`
	ID               string `json:"id"`
	CreatedAt        int64  `json:"created_at"`
	PreserveHost     bool   `json:"preserve_host"`
	Name             string `json:"name"`
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

func (s *ApisService) GetAll(opt *ApisGetAllOptions) (*Apis, *http.Response, error) {

}