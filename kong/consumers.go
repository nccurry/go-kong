package kong

import (
	"errors"
	"fmt"
	"net/http"
)

type ConsumersService service

type Consumers struct {
	Data  []Consumer `json:"consumer,omitempty"`
	Total int        `json:"total,omitempty"`
	Next  string     `json:"next,omitempty"`
}

type Consumer struct {
	ID        string `json:"id,omitempty"`
	Username  string `json:"username,omitempty"`
	CustomID  string `json:"custom_id,omitempty"`
	CreatedAt int    `json:"created_at,omitempty"`
}

// Get returns a single Kong consumer. The consumer username or id can be used
func (s *ConsumersService) Get(consumer string) (*Consumer, *http.Response, error) {
	u := fmt.Sprintf("consumers/%v", consumer)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	uResp := new(Consumer)
	resp, err := s.client.Do(req, uResp)
	if err != nil {
		return nil, resp, err
	}

	return uResp, resp, err
}

func (s *ConsumersService) Patch(consumer *Consumer) (*http.Response, error) {
	var u string
	if consumer.ID != "" {
		u = fmt.Sprintf("consumers/%v", consumer.ID)
	} else if consumer.Username != "" {
		u = fmt.Sprintf("consumers/%v", consumer.Username)
	} else {
		return nil, errors.New("You must specify either consumer username or id.")
	}

	req, err := s.client.NewRequest("PATCH", u, consumer)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)

	return resp, err
}

// Delete removes a Kong consumer. The username or id field can be used for consumer
func (s *ConsumersService) Delete(consumer string) (*http.Response, error) {
	u := fmt.Sprintf("apis/%v", consumer)

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

func (s *ConsumersService) Post(consumer *Consumer) (*http.Response, error) {
	req, err := s.client.NewRequest("POST", "consumers", consumer)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)

	return resp, err
}

type ConsumersGetAllOptions struct {
	ID       string `url:"id,omitempty"`        // A filter on the list based on the consumer id field.
	CustomID string `url:"custom_id,omitempty"` // A filter on the list based on the consumer custom_id field.
	Username string `url:"username,omitempty"`  // A filter on the list based on the consumer username field.
	Size     int    `url:"size,omitempty"`      // A limit on the number of objects to be returned.
	Offset   string `url:"offset,omitempty"`    // A cursor used for pagination. offset is an object identifier that defines a place in the list.
}

func (s *ConsumersService) GetAll(opt *ConsumersGetAllOptions) (*Consumers, *http.Response, error) {
	u, err := addOptions("consumers", opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	consumers := new(Consumers)
	resp, err := s.client.Do(req, consumers)
	if err != nil {
		return nil, resp, err
	}

	return consumers, resp, err
}
