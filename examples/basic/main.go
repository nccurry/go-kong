package main

import (
	"github.com/spicyusername/go-kong/kong"
	"log"
)

func main() {

	client, err := kong.NewClient(nil, "http://localhost:8001")
	if err != nil {
		log.Fatal(err)
	}

	/*
	api := &kong.Api{
		Name:"Test3",
		RequestPath:"/test3",
		UpstreamURL:"http://test.com:8080",
	}
	_, err = client.Apis.Post(api)
	if err != nil {
		log.Fatal(err)
	}
*/
	opt := &kong.ApisGetAllOptions{UpstreamURL: "http://test.com:8080"}
	apis, _, err := client.Apis.GetAll(opt)
	if err != nil {
		log.Fatal(err)
	}

	aclPlugin := &kong.ACLPlugin{
		Config: kong.ACLConfig{Whitelist: []string{"admins"}},
		Plugin: kong.Plugin{Name: "acl"},
	}
	plugin := aclPlugin.ToPlugin()
	_, err = client.Plugins.Post(plugin)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v", apis)
}