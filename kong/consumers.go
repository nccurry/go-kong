package kong

import (
	"errors"
	"fmt"
	"net/http"
)

// ConsumersService handles communication with Kong's '/consumers' resource.
type ConsumersService service

// Consumers represents the object returned from Kong when querying for multiple consumer objects.
//
// In cases where the number of objects returned exceeds the maximum, Next holds the
// URI for the next set of results.
// i.e. "http://localhost:8001/consumers/?size=2&offset=4d924084-1adb-40a5-c042-63b19db421d1"
type Consumers struct {
	Data  []Consumer `json:"consumer,omitempty"`
	Total int        `json:"total,omitempty"`
	Next  string     `json:"next,omitempty"`
}

// Consumer represents a single Kong consumer object
type Consumer struct {
	ID        string `json:"id,omitempty"`
	Username  string `json:"username,omitempty"`
	CustomID  string `json:"custom_id,omitempty"`
	CreatedAt int    `json:"created_at,omitempty"`
}

// ConsumersService.Get queries for a single Kong consumer object, by name or id.
//
// Equivalent to GET /consumers/{name or id}
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

// ConsumersService.Patch updates an existing Kong consumer object.
// At least one of consumer.Username or consumer.ID must be specified
// in the passed *Consumer parameter.
//
// Equivalent to PATCH /consumers/{username or id}
func (s *ConsumersService) Patch(consumer *Consumer) (*http.Response, error) {
	var u string
	if consumer.ID != "" {
		u = fmt.Sprintf("consumers/%v", consumer.ID)
	} else if consumer.Username != "" {
		u = fmt.Sprintf("consumers/%v", consumer.Username)
	} else {
		return nil, errors.New("At least one of consumer.Username or consumer.ID must be specified")
	}

	req, err := s.client.NewRequest("PATCH", u, consumer)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)

	return resp, err
}

// ConsumersService.Delete deletes a single Kong consumer object, by name or id.
//
// Equivalent to DELETE /consumers/{username or id}
func (s *ConsumersService) Delete(consumer string) (*http.Response, error) {
	u := fmt.Sprintf("consumers/%v", consumer)

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

// ConsumersService.Post creates a new Kong consumer object.
//
// Equivalent to POST /consumers
func (s *ConsumersService) Post(consumer *Consumer) (*http.Response, error) {
	req, err := s.client.NewRequest("POST", "consumers", consumer)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)

	return resp, err
}

// ConsumersGetAllOptions specifies optional filter parameters to the
// ConsumersService.GetAll method.
//
// Additional information about filtering options can be found in the
// Kong documentation at:
// https://getkong.org/docs/0.9.x/admin-api/#list-consumers
type ConsumersGetAllOptions struct {
	ID       string `url:"id,omitempty"`        // A filter on the list based on the consumer id field.
	CustomID string `url:"custom_id,omitempty"` // A filter on the list based on the consumer custom_id field.
	Username string `url:"username,omitempty"`  // A filter on the list based on the consumer username field.
	Size     int    `url:"size,omitempty"`      // A limit on the number of objects to be returned.
	Offset   string `url:"offset,omitempty"`    // A cursor used for pagination. offset is an object identifier that defines a place in the list.
}

// ConsumersService.GetAll queries for all Kong consumer objects.
// This query can be filtered by supplying the ConsumersGetAllOptions struct.
//
// Equivalent to GET /consumers?uri=params&from=opt
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

type ACLConsumerConfig struct {
	Group string `json:"group"`
}

func (s *ConsumersService) ConfigurePlugin(consumer, plugin string, config interface{}) (*http.Response, error){
	u := fmt.Sprintf("consumers/%v/%v", consumer, plugin)

	req, err := s.client.NewRequest("POST", u, config)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)

	return resp, err
}