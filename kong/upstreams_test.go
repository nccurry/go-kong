package kong

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestUpstream_Get(t *testing.T) {
	stubSetup()
	defer stubTeardown()

	mux.HandleFunc("/upstreams/i", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":"i"}`)
	})

	upstream, _, err := client.Upstreams.Get("i")
	if err != nil {
		t.Errorf("Upstreams.Get returned error: %v", err)
	}

	want := &Upstream{ID: "i"}
	if !reflect.DeepEqual(upstream, want) {
		t.Errorf("Upstreams.Get returned %+v, want %+v", upstream, want)
	}
}

func TestUpstream_Get_invalidUpstream(t *testing.T) {
	stubSetup()
	defer stubTeardown()

	_, _, err := client.Upstreams.Get("%")
	testURLParseError(t, err)
}

func TestUpstream_Get_badStatusCode(t *testing.T) {
	stubSetup()
	defer stubTeardown()

	mux.HandleFunc("/upstreams/i", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		fmt.Fprint(w, `{"error":"e"}`)
	})

	_, _, err := client.Upstreams.Get("i")
	if err == nil {
		t.Error("Expected error to be returned")
	}
}

func TestUpstreams_Patch_byName(t *testing.T) {
	stubSetup()
	defer stubTeardown()

	input := sampleUpstream()

	mux.HandleFunc("/upstreams/"+input.Name, func(w http.ResponseWriter, r *http.Request) {
		v := new(Upstream)
		json.NewDecoder(r.Body).Decode(v)
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		testMethod(t, r, "PATCH")
	})

	_, err := client.Upstreams.Patch(input)
	if err != nil {
		t.Errorf("Upstreams.Patch returned error: %v", err)
	}
}

func TestUpstreams_Patch_byID(t *testing.T) {
	stubSetup()
	defer stubTeardown()

	input := sampleUpstream()

	mux.HandleFunc("/upstreams/"+input.Name, func(w http.ResponseWriter, r *http.Request) {
		v := new(Upstream)
		json.NewDecoder(r.Body).Decode(v)
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		testMethod(t, r, "PATCH")
	})

	_, err := client.Upstreams.Patch(input)
	if err != nil {
		t.Errorf("Upstreams.Patch returned error: %v", err)
	}
}

func TestUpstream_Patch_invalidUpstream(t *testing.T) {
	stubSetup()
	defer stubTeardown()

	input := &Upstream{Name: "%"}
	_, err := client.Upstreams.Patch(input)
	testURLParseError(t, err)
}

func TestUpstream_Patch_missingIDOrName(t *testing.T) {
	stubSetup()
	defer stubTeardown()

	input := &Upstream{Slots: 100}
	_, err := client.Upstreams.Patch(input)
	if err == nil {
		t.Error("Expected error to be returned")
	}
}

func TestUpstream_Patch_badStatusCode(t *testing.T) {
	stubSetup()
	defer stubTeardown()

	mux.HandleFunc("/upstreams/i", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		fmt.Fprint(w, `{"error":"e"}`)
	})

	input := sampleUpstream()

	_, err := client.Upstreams.Patch(input)
	if err == nil {
		t.Error("Expected error to be returned")
	}
}

func TestUpstream_Delete(t *testing.T) {
	stubSetup()
	defer stubTeardown()

	mux.HandleFunc("/upstreams/i", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	_, err := client.Upstreams.Delete("i")
	if err != nil {
		t.Errorf("Upstreams.Delete returned error: %v", err)
	}
}

func TestUpstream_Delete_invalidUpstream(t *testing.T) {
	_, err := client.Upstreams.Delete("%")
	testURLParseError(t, err)
}

func TestUpstream_Delete_badStatusCode(t *testing.T) {
	stubSetup()
	defer stubTeardown()

	mux.HandleFunc("/upstreams/i", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		fmt.Fprint(w, `{"error":"e"}`)
	})

	_, err := client.Upstreams.Delete("i")
	if err == nil {
		t.Error("Expected error to be returned")
	}
}

func TestUpstream_Post(t *testing.T) {
	stubSetup()
	defer stubTeardown()

	input := sampleUpstream()

	mux.HandleFunc("/upstreams", func(w http.ResponseWriter, r *http.Request) {
		v := new(Upstream)
		json.NewDecoder(r.Body).Decode(v)
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		testMethod(t, r, "POST")
	})

	_, err := client.Upstreams.Post(input)
	if err != nil {
		t.Errorf("Upstreams.Post returned error: %v", err)
	}
}

func TestUpstream_Post_badStatusCode(t *testing.T) {
	stubSetup()
	defer stubTeardown()

	mux.HandleFunc("/upstreams", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		fmt.Fprint(w, `{"error":"e"}`)
	})

	input := sampleUpstream()

	_, err := client.Upstreams.Post(input)
	if err == nil {
		t.Error("Expected error to be returned")
	}
}

func sampleUpstream() *Upstream {
	return &Upstream{
		Name:      "upstreamName",
		ID:        "upstreamID",
		Slots:     4,
		Orderlist: []int{4, 1, 3, 2},
	}
}
