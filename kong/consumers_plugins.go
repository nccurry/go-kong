package kong

import (
	"fmt"
	"net/http"
)

type ConsumersPlugins struct {
	ACL *ConsumersACLService
	JWT *ConsumersJWTService
}

type ConsumersACLService service

type ConsumerACLConfigs struct {
	Data  []ConsumerACLConfig `json:"data,omitempty"`
	Total int                 `json:"total,omitempty"`
}

type ConsumerACLConfig struct {
	ConsumerID string `json:"consumer_id,omitempty"`
	CreatedAt  int    `json:"created_at,omitempty"`
	Group      string `json:"group,omitempty"`
	ID         string `json:"id,omitempty"`
}

func (s *ConsumersACLService) Post(consumer string, config *ConsumerACLConfig) (*http.Response, error) {
	u := fmt.Sprintf("consumers/%v/acls", consumer)

	req, err := s.client.NewRequest("POST", u, config)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)

	return resp, err
}

func (s *ConsumersACLService) Get(consumer string) (*ConsumerACLConfigs, *http.Response, error) {
	u := fmt.Sprintf("consumers/%v/acls", consumer)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	uResp := new(ConsumerACLConfigs)
	resp, err := s.client.Do(req, uResp)
	if err != nil {
		return nil, resp, err
	}

	return uResp, resp, err
}

func (s *ConsumersACLService) Delete(consumer, id string) (*http.Response, error) {
	u := fmt.Sprintf("consumers/%v/acls/%v", consumer, id)

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

type ConsumersJWTService service

type ConsumerJWTConfigs struct {
	Data  []ConsumerJWTConfig `json:"data,omitempty"`
	Total int                 `json:"total,omitemtpy"`
}

type ConsumerJWTConfig struct {
	Key          string `json:"key,omitempty"`
	Algorithm    string `json:"algorithm,omitempty"`
	RSAPublicKey string `json:"rsa_public_key,omitempty"`
	Secret       string `json:"secret,omitempty"`
}

func (s *ConsumersJWTService) Post(consumer string, config *ConsumerJWTConfig) (*http.Response, error) {
	u := fmt.Sprintf("consumers/%v/jwt", consumer)

	req, err := s.client.NewRequest("POST", u, config)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)

	return resp, err
}

func (s *ConsumersJWTService) Get(consumer string) (*ConsumerJWTConfigs, *http.Response, error) {
	u := fmt.Sprintf("consumers/%v/jwt", consumer)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	uResp := new(ConsumerJWTConfigs)
	resp, err := s.client.Do(req, uResp)
	if err != nil {
		return nil, resp, err
	}

	return uResp, resp, err
}

func (s *ConsumersJWTService) Delete(consumer, id string) (*http.Response, error) {
	u := fmt.Sprintf("consumers/%v/jwt/%v", consumer, id)

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
