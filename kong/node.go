package kong

import (
	"net/http"
)

type NodeService service

type Node struct {
	Configuration map[string]interface{} `json:"configuration,omitempty"`
	Hostname      string                 `json:"hostname,omitempty"`
	LuaVersion    string                 `json:"lua_version,omitempty"`
	Plugins       struct {
		AvailableOnServer map[string]bool `json:"available_on_server,omitempty"`
		EnabledInCluster  map[string]bool `json:"enabled_in_cluster,omitempty"`
	} `json:"plugins,omitempty"`
	PRNGSeeds map[string]int `json:"prng_seeds,omitempty"`
	Tagline   string         `json:"tagline,omitempty"`
	Timers    map[string]int `json:"timers,omitempty"`
	Version   string         `json:"version,omitempty"`
}

type Status struct {
	Database map[string]int `json:"database,omitempty"`
	Server   map[string]int `json:"server,omitempty"`
}

func (s *NodeService) Get() (*Node, *http.Response, error) {
	req, err := s.client.NewRequest("GET", "", nil)
	if err != nil {
		return nil, nil, err
	}

	uResp := new(Node)
	resp, err := s.client.Do(req, uResp)
	if err != nil {
		return nil, resp, err
	}

	return uResp, resp, err
}

func (s *NodeService) GetStatus() (*Status, *http.Response, error) {
	req, err := s.client.NewRequest("GET", "status", nil)
	if err != nil {
		return nil, nil, err
	}

	uResp := new(Status)
	resp, err := s.client.Do(req, uResp)
	if err != nil {
		return nil, resp, err
	}

	return uResp, resp, err
}

type ClusterService service

type Cluster struct {
	Total int             `json:"total,omitemtpy"`
	Data  []ClusterMember `json:"data,omitempty"`
}

type ClusterMember struct {
	Address string `json:"address,omitempty"`
	Name    string `json:"name,omitempty"`
	Status  string `json:"status,omitempty"`
}

func (s *ClusterService) Get() (*Cluster, *http.Response, error) {
	req, err := s.client.NewRequest("GET", "cluster", nil)
	if err != nil {
		return nil, nil, err
	}

	uResp := new(Cluster)
	resp, err := s.client.Do(req, uResp)
	if err != nil {
		return nil, resp, err
	}

	return uResp, resp, err
}

func (s *ClusterService) Delete(clusterMember *ClusterMember) (*http.Response, error) {
	req, err := s.client.NewRequest("DELETE", "cluster", clusterMember)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, err
}
