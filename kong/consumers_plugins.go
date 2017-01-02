package kong

import (
	"fmt"
	"net/http"
)

type ConsumersACLService service

type ACLConsumerConfigs struct {
	Data  []ACLConsumerConfig `json:"data,omitempty"`
	Total int                 `json:"total,omitempty"`
}

type ACLConsumerConfig struct {
	ConsumerID string `json:"consumer_id,omitempty"`
	CreatedAt  int    `json:"created_at,omitempty"`
	Group      string `json:"group,omitempty"`
	ID         string `json:"id,omitempty"`
}

func (s *ConsumersACLService) Configure(consumer string, config *ACLConsumerConfig) (*http.Response, error) {
	u := fmt.Sprintf("consumers/%v/acls", consumer)

	req, err := s.client.NewRequest("POST", u, config)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)

	return resp, err
}

func (s *ConsumersACLService) Get(consumer string) (*ACLConsumerConfigs, *http.Response, error) {
	u := fmt.Sprintf("consumers/%v/acls", consumer)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	uResp := new(ACLConsumerConfigs)
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
