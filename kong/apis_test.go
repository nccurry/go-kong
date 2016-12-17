package kong

import "testing"

func TestAPI_marshal(t *testing.T) {
	testJSONMarshal(t, &Api{}, "{}")

	a := &Api{
		Name: "n",
		CreatedAt: 1,
		ID: "i",
		PreserveHost: true,
		RequestPath: "r",
		StripRequestPath: true,
		UpstreamURL: "u",
	}
	want := `{
		"name": "n",
		"created_at": 1,
		"id": "i",
		"preserve_host": true,
		"request_path": "r",
		"strip_request_path": true,
		"upstream_url": "u"
	}`
	testJSONMarshal(t, a, want)
}
