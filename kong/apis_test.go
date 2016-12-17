package kong

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

/*
func TestAPI_marshal(t *testing.T) {
	testJSONMarshal(t, &Api{}, "{}")

	a := &Api{
		Name:             "n",
		CreatedAt:        1,
		ID:               "i",
		PreserveHost:     true,
		RequestPath:      "r",
		StripRequestPath: true,
		UpstreamURL:      "u",
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
*/
func TestApisService_Get(t *testing.T) {
	stubSetup()
	defer stubTeardown()

	mux.HandleFunc("/apis/i", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":"i"}`)
	})

	user, _, err := client.Apis.Get("i")
	if err != nil {
		t.Errorf("Users.Get returned error: %v", err)
	}

	want := &Api{ID: "i"}
	if !reflect.DeepEqual(user, want) {
		t.Errorf("Users.Get returned %+v, want %+v", user, want)
	}
}

func TestUsersService_Get_invalidUser(t *testing.T) {
	_, _, err := client.Apis.Get("%")
	testURLParseError(t, err)
}

func TestApisService_Patch(t *testing.T) {
	stubSetup()
	defer stubTeardown()

	input := &Api{ID: "i"}

	mux.HandleFunc("/apis/i", func(w http.ResponseWriter, r *http.Request) {
		v := new(Api)
		json.NewDecoder(r.Body).Decode(v)
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		testMethod(t, r, "PATCH")

	})
}
