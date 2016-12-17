package main

import (
	"go-kong/kong"
	"log"
)

func createNewAPI(client *kong.Client) {
	api := &kong.Api{
		Name:"Test3",
		RequestPath:"/test3",
		UpstreamURL:"http://test.com:8080",
	}
	_, err := client.Apis.Post(api)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	client, err := kong.NewClient(nil, "http://localhost:8001")
	if err != nil {
		log.Fatal(err)
	}

	//createNewAPI(client)
	aclConfig := &kong.ACLConfig{Whitelist: []string{"admins"}}
	aclPlugin := &kong.ACLPlugin{
		Config: aclConfig,
		Plugin: kong.Plugin{Name: "acl"},
	}
	plugin := aclPlugin.ToPlugin()

	_, err = client.Plugins.Post(plugin)
	if err != nil {
		log.Fatal(err)
	}

}