package kong

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"os"
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

	api, _, err := client.Apis.Get("i")
	if err != nil {
		t.Errorf("Apis.Get returned error: %v", err)
	}

	want := &Api{ID: "i"}
	if !reflect.DeepEqual(api, want) {
		t.Errorf("Apis.Get returned %+v, want %+v", api, want)
	}
}

func TestApisService_Get_invalidApi(t *testing.T) {
	_, _, err := client.Apis.Get("%")
	testURLParseError(t, err)
}

func TestApisService_Get_badStatusCode(t *testing.T) {
	stubSetup()
	defer stubTeardown()

	mux.HandleFunc("/apis/i", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		fmt.Fprint(w, `{"error":"e"}`)
	})

	_, _, err := client.Apis.Get("i")
	if err == nil {
		t.Error("Expected error to be returned")
	}
}

func TestApisService_Patch_byName(t *testing.T) {
	stubSetup()
	defer stubTeardown()

	input := &ApiRequest{Name: "n"}

	mux.HandleFunc("/apis/n", func(w http.ResponseWriter, r *http.Request) {
		v := new(ApiRequest)
		json.NewDecoder(r.Body).Decode(v)
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		testMethod(t, r, "PATCH")

	})

	_, err := client.Apis.Patch(input)
	if err != nil {
		t.Errorf("Apis.Patch returned error: %v", err)
	}
}

func TestApisService_Patch_byID(t *testing.T) {
	stubSetup()
	defer stubTeardown()

	input := &ApiRequest{ID: "i"}

	mux.HandleFunc("/apis/i", func(w http.ResponseWriter, r *http.Request) {
		v := new(ApiRequest)
		json.NewDecoder(r.Body).Decode(v)
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		testMethod(t, r, "PATCH")

	})

	_, err := client.Apis.Patch(input)
	if err != nil {
		t.Errorf("Apis.Patch returned error: %v", err)
	}
}

func TestApisService_Patch_invalidApi(t *testing.T) {
	input := &ApiRequest{Name: "%"}
	_, err := client.Apis.Patch(input)
	testURLParseError(t, err)
}

func TestApisService_Patch_missingIDOrName(t *testing.T) {
	input := &ApiRequest{RequestPath: "r"}
	_, err := client.Apis.Patch(input)
	if err == nil {
		t.Error("Expected error to be returned")
	}
}

func TestApisService_Patch_badStatusCode(t *testing.T) {
	stubSetup()
	defer stubTeardown()

	mux.HandleFunc("/apis/i", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		fmt.Fprint(w, `{"error":"e"}`)
	})

	input := &ApiRequest{ID: "i"}

	_, err := client.Apis.Patch(input)
	if err == nil {
		t.Error("Expected error to be returned")
	}
}

func TestApisService_Delete(t *testing.T) {
	stubSetup()
	defer stubTeardown()

	mux.HandleFunc("/apis/i", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	_, err := client.Apis.Delete("i")
	if err != nil {
		t.Errorf("Apis.Delete returned error: %v", err)
	}
}

func TestApisService_Delete_invalidApi(t *testing.T) {
	_, err := client.Apis.Delete("%")
	testURLParseError(t, err)
}

func TestApisService_Delete_badStatusCode(t *testing.T) {
	stubSetup()
	defer stubTeardown()

	mux.HandleFunc("/apis/i", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		fmt.Fprint(w, `{"error":"e"}`)
	})

	_, err := client.Apis.Delete("i")
	if err == nil {
		t.Error("Expected error to be returned")
	}
}

func TestApisService_Post(t *testing.T) {
	stubSetup()
	defer stubTeardown()

	input := &ApiRequest{ID: "i"}

	mux.HandleFunc("/apis", func(w http.ResponseWriter, r *http.Request) {
		v := new(ApiRequest)
		json.NewDecoder(r.Body).Decode(v)
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		testMethod(t, r, "POST")

	})

	_, err := client.Apis.Post(input)
	if err != nil {
		t.Errorf("Apis.Post returned error: %v", err)
	}
}

func TestApisService_Post_badStatusCode(t *testing.T) {
	stubSetup()
	defer stubTeardown()

	mux.HandleFunc("/apis", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		fmt.Fprint(w, `{"error":"e"}`)
	})

	input := &ApiRequest{ID: "i"}

	_, err := client.Apis.Post(input)
	if err == nil {
		t.Error("Expected error to be returned")
	}
}

func TestApisService_GetAll(t *testing.T) {
	stubSetup()
	defer stubTeardown()

	v := &Apis{Total: 1, Next: "n", Data: []*Api{{ID: "i"}}}

	mux.HandleFunc("/apis", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"offset": "o", "request_host": "r"})
		json.NewEncoder(w).Encode(v)
	})

	opt := &ApisGetAllOptions{Offset: "o", RequestHost: "r"}
	apis, _, err := client.Apis.GetAll(opt)
	if err != nil {
		t.Errorf("Apis.GetAll returned error: %v", err)
	}

	want := &Apis{Total: 1, Next: "n", Data: []*Api{{ID: "i"}}}
	if !reflect.DeepEqual(apis, want) {
		t.Errorf("Apis.GetAll returned %+v, want %+v", apis, want)
	}
}

func TestApisService_GetAll_badStatusCode(t *testing.T) {
	stubSetup()
	defer stubTeardown()

	mux.HandleFunc("/apis", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		fmt.Fprint(w, `{"error":"e"}`)
	})

	_, _, err := client.Apis.GetAll(nil)
	if err == nil {
		t.Error("Expected error to be returned")
	}
}

func TestGet(t *testing.T) {
	id := "8dcf8d7a-f173-4367-949f-253320f3eedf"
 	baseUrl := os.Getenv("KONG_ADMIN_URL")
	c, err := NewClient(&http.Client{}, baseUrl)
	fmt.Errorf("err: %v",err)
	api, _, err := c.Apis.Get(id)
	fmt.Errorf("api : %v", api)
}


func TestPost(t *testing.T) {
	baseUrl := os.Getenv("KONG_ADMIN_URL")
	c, err := NewClient(&http.Client{}, baseUrl)
	fmt.Errorf("err: %v",err)
	apiRequest := &ApiRequest{
		Name:                   "Example",
		Hosts:                  []string{"example.com"},
		Uris:                   []string{"/example"},
		UpstreamURL:            "http://localhost:4140/testservice",
		StripUri:               true,
		PreserveHost:           true,
		Retries:                3,
		UpstreamConnectTimeout: 1000,
		UpstreamSendTimeout:    2000,
		UpstreamReadTimeout:    3000,
		HttpsOnly:              true,
		HttpIfTerminated:       true,
	}
	_, err = c.Apis.Post(apiRequest)
	fmt.Errorf("api : %v", err)
}
