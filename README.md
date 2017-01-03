# go-kong #

[![GoDoc](https://godoc.org/github.com/spicyusername/go-kong/kong?status.svg)](https://godoc.org/github.com/spicyusername/go-kong/kong) 
[![Build Status](https://travis-ci.org/spicyusername/go-kong.svg?branch=master)](https://travis-ci.org/spicyusername/go-kong) 
[![Coverage Status](https://coveralls.io/repos/github/spicyusername/go-kong/badge.svg?branch=master)](https://coveralls.io/github/spicyusername/go-kong?branch=master) 
[![Go Report Card](https://goreportcard.com/badge/github.com/spicyusername/go-kong)](https://goreportcard.com/report/github.com/spicyusername/go-kong) 

go-kong is a Go client library for accessing Mashape's [Kong API](https://getkong.org/docs/0.9.x/admin-api/).

## Table of Contents ##

* [Installation](#installation)
* [Usage](#usage)
* [Client Objects](#client-objects)
    * [Node](#node)
    * [Cluster](*cluster)
    * [Apis](#apis)  
    * [Consumers](#consumers)
    * [Plugins](#plugins)
    * [Consumers Plugins](#consumers-plugins)
* [Handling Errors](#handling-errors)
* [Filtering with Query Parameters](#filtering-with-query-parameters)
* [Working with Plugin Definitions](#working-with-plugin-definitions)
* [To-Do](#to-do)

## Installation ##

```bash
go get "github.com/spicyusername/go-kong/kong"
```

## Usage ##

Construct a new Kong client, then use the various services on the client to
access different parts of the Kong API. 

For example:

```go
package main

import (
        "log"
        "github.com/spicyusername/go-kong/kong"
)

func main() {
    // Create new client
    client, _ := kong.NewClient(nil, "http://localhost:8001/")
    
    // Get information about the 'backend' api
    api, _, _ := client.Apis.Get("backend")
    log.Println(api.ID)

    // Create a new api called 'middletier'
    mtApi := &kong.Api{Name: "middletier", RequestPath: "/mt/v0", UpstreamURL: "http://mt.my.org:8080"}
    _, err := client.Apis.Post(mtApi)
    
    // Handle 409 error separately
    if ok := err.(kong.ConflictError); ok {
        log.Printf("Endpoint with name %s already exists.", mtApi.Name)
    } else if err != nil {
        log.Fatal(err)
    }
    
    // Get all consumer objects
    consumers, _, _ := client.Consumers.GetAll(nil)
    for _, v := range consumers.Data {
        log.Println(v.Username)
    }

    // Apply ACL plugin to all all apis
    aclConfig := &kong.ACLConfig{Whitelist: []string{"users", "admins"}, Blacklist: []string{"blocked"}}
    plugin := &kong.Plugin{Name: "acl", Config: kong.ToMap(aclConfig)}
    _, err = client.Plugins.Post(plugin)
    if err != nil {
        log.Fatal(err)
    }
    
    // Add ACL group to consumer
    aclConsumerConfig := &kong.ConsumerACLConfig{Group: "users"}
    consumerName := "paul.atreides"
    _, err = client.Consumers.ACL.Configure(consumerName, aclConsumerConfig)
    
    // Handle 404 separately
    if ok := err.(kong.NotFoundError); ok {
        log.Printf("Could not find consumer with name %s", consumerName)
    } else if err != nil {
        log.Fatal(err)
    }
}
```

## Client Objects ##

#### Node ####

```go
// GET /
node, resp, err := client.Node.Get()

// GET /status
status, resp, err := client.Node.GetStatus()
```

```go
type Node struct {
	Configuration map[string]interface{} `json:"configuration,omitempty"`
	Hostname string `json:"hostname,omitempty"`
	LuaVersion string `json:"lua_version,omitempty"`
	Plugins struct{
		AvailableOnServer map[string]bool `json:"available_on_server,omitempty"`
		EnabledInCluster map[string]bool `json:"enabled_in_cluster,omitempty"`
	} `json:"plugins,omitempty"`
	PRNGSeeds map[string]int `json:"prng_seeds,omitempty"`
	Tagline string `json:"tagline,omitempty"`
	Timers map[string]int `json:"timers,omitempty"`
	Version string `json:"version,omitempty"`
}

type Status struct {
	Database map[string]int `json:"database,omitempty"`
	Server map[string]int `json:"server,omitempty"`
}
```

#### Cluster ####

```go
// GET /cluster
cluster, resp, err := client.Cluster.Get()

// DELETE /cluster
cluster := &kong.Cluster{Name: "clusternode01"}
resp, err := client.Cluster.Delete(cluster)
```

```go
type Cluster struct {
	Total int `json:"total,omitemtpy"`
	Data []ClusterMember `json:"data,omitempty"`
}

type ClusterMember struct {
	Address string `json:"address,omitempty"`
	Name string `json:"name,omitempty"`
	Status string `json:"status,omitempty"`
}
```

#### Apis ####

```go
// GET /apis
apis, resp, err := client.Apis.GetAll(nil)

// GET /apis?size=10&mapKey=myapi
opt := &ApisGetAllOptions{Size: 10, Name: "myapi"}
apis, resp, err := client.Apis.GetAll(opt)

// GET /apis/myapi
api, resp, err := client.Apis.Get("myapi")

// POST /apis
api := &Api{Name: "myapi", RequestPath: "/myapi", UpstreamURL: "http:myapi:8080"}
resp, err := client.Apis.Post(api)

// PATCH /apis/myapi
api := &Api{Name: "myapi", RequestPath: "/myapi", UpstreamURL: "http:myapi:8080"}
resp, err := client.Apis.Patch(api)

// DELETE /apis/myapi
resp, err := client.Apis.Delete("myapi")
```

```go
type Apis struct {
	Data   []Api  `json:"data,omitempty"`
	Total  int    `json:"total,omitempty"`
	Next   string `json:"next,omitempty"`
	Offset string `json:"offset,omitempty"`
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

type ApisGetAllOptions struct {
	ID          string `url:"id,omitempty"`           
	Name        string `url:"name,omitempty"`         
	RequestHost string `url:"request_host,omitempty"` 
	RequestPath string `url:"request_path,omitempty"`
	UpstreamURL string `url:"upstream_url,omitempty"`
	Size        int    `url:"size,omitempty"`        
	Offset      string `url:"offset,omitempty"`      
}
```

#### Consumers ####

```go
// GET /consumers
consumers, resp, err := client.Consumers.GetAll(nil)

// GET /consumers?size=10&custom_id=nickname
opt := &ConsumersGetAllOptions{Size: 10, CustomID: "nickname"}
consumers, resp, err := client.Consumers.GetAll(opt)

// GET /consumers/admin
consumer, resp, err := client.Consumers.Get("admin")

// POST /consumers
consumer := &Consumer{Username: "admin"}
resp, err := client.Consumers.Post(consumer)

// PATCH /consumers/admin
consumer := &Consumer{CustomID: "superuser"}
resp, err := client.Consumers.Patch(consumer)

// DELETE /consumers/admin
resp, err := client.Consumers.Delete("admin")
```

```go
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

type ConsumersGetAllOptions struct {
	ID       string `url:"id,omitempty"`       
	CustomID string `url:"custom_id,omitempty"`
	Username string `url:"username,omitempty"`  
	Size     int    `url:"size,omitempty"`   
	Offset   string `url:"offset,omitempty"`
}
```

#### Plugins ####

Because of the generic nature of Kong's plugin object, the Config field of
the Plugin struct is defined as follows:
```go
type Plugin struct {
    ...
    Config map[string]interface{} `json:"config,omitempty"`
    ...
}
```

All of the various Kong plugin configurations have (eventually)
been defined. Two helper functions have been written to help marshal plugin configurations to/from 
more specific configuration structs to the more generic ```map[string]interface{}```

For example:
```go
// Convert struct to map[string]interface{} expected by kong.Plugin
aclConfig := &kong.ACLConfig{Whitelist: []string{"users", "admins"}, Blacklist: []string{"blocked"}}
plugin := &kong.Plugin{Name: "acl", Config: kong.ToMap(aclConfig)}

// Convert map[string]interface{} returned by client.Plugins.Get to struct
plugin, _, _ := client.Plugins.Get(id)
aclConfig := new(ACLConfig)
err := kong.FromMap(aclConfig, plugin.Config)
```



```go
// GET /plugins
plugins, resp, err := client.Plugins.GetAll(nil)

// GET /plugins?size=15&consumer_id=4def15f5-0697-4956-a2b0-9ae079b686bb
opt := &kong.PluginsGetAllOptions{Size: 15, ConsumerID: "4def15f5-0697-4956-a2b0-9ae079b686bb"}
plugins, resp, err := client.Plugins.GetAll(opt)

// GET /plugins/4def15f5-0697-4956-a2b0-9ae079b686bb
plugin, resp, err := client.Plugins.Get("4def15f5-0697-4956-a2b0-9ae079b686bb")

// GET /plugins/enabled
enabled, resp, err := client.Plugins.GetEnabled()

// GET /plugins/schema/acl
schema, resp, err := client.Plugins.GetSchema("acl")

// POST /plugins
aclConfig := &kong.ACLConfig{Whitelist: []string{"users", "admins"}, Blacklist: []string{"blocked"}}
plugin := &kong.Plugin{Name: "acl", Config: kong.ToMap(aclConfig)}
resp, err := client.Plugins.Post(plugin)

// PATCH /plugins
aclConfig := &kong.ACLConfig{Whitelist: []string{"users", "admins"}, Blacklist: []string{"blocked"}}
plugin := &kong.Plugin{Name: "acl", Config: kong.ToMap(aclConfig)}
resp, err := client.Plugins.Patch(plugin)

// DELETE /plugins/4def15f5-0697-4956-a2b0-9ae079b686bb
resp, err := client.Plugins.Delete("4def15f5-0697-4956-a2b0-9ae079b686bb")
```

```go
type Plugins struct {
	Data  []Plugin `json:"data,omitempty"`
	Total int      `json:"total,omitempty"`
	Next  string   `json:"next,omitempty"`
}

type Plugin struct {
	ID         string                 `json:"id,omitempty"`
	Name       string                 `json:"name,omitempty"`
	CreatedAt  int                    `json:"created_at,omitempty"`
	Enabled    bool                   `json:"enabled,omitempty"`
	ApiID      string                 `json:"api_id,omitempty"`
	ConsumerID string                 `json:"consumer_id,omitempty"`
	Config     map[string]interface{} `json:"config,omitempty"`
}

type PluginsGetAllOptions struct {
	ID         string `url:"id,omitempty"`          // A filter on the list based on the id field.
	Name       string `url:"name,omitempty"`        // A filter on the list based on the name field.
	ApiID      string `url:"api_id,omitempty"`      // A filter on the list based on the api_id field.
	ConsumerID string `url:"consumer_id,omitempty"` // A filter on the list based on the consumer_id field.
	Size       int    `url:"size,omitempty"`        // A limit on the number of objects to be returned.
	Offset     string `url:"offset,omitempty"`      // A cursor used for pagination. offset is an object identifier that defines a place in the list.

}
```
In addition to the generic plugin struct definitions, there are many more structures defined for each plugin
configuration in [plugins.go](kong/plugins.go)

#### Consumers Plugins ####

```go
```

## Handling Errors ##

## Filtering with Query Parameters ##

## Working with Plugin Definitions ##

## To-Do ##
* Finish the README.md
* Fuller Unit-testing
* Represent all plugin object configs via structs
* Represent all consumer plugin configs via structs
