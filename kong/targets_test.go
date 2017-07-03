package kong

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

var upstreamName = "testupstream"

func TestTargets_Delete(t *testing.T) {
	stubSetup()
	defer stubTeardown()

	mux.HandleFunc(fmt.Sprintf("/upstreams/%s/targets/i", upstreamName), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	_, err := client.Targets.Delete(upstreamName, "i")
	if err != nil {
		t.Errorf("Targets.Delete returned error: %v", err)
	}
}

func TestTargets_Delete_invalidTarget(t *testing.T) {
	stubSetup()
	defer stubTeardown()

	_, err := client.Targets.Delete(upstreamName, "%")
	testURLParseError(t, err)
}

func TestTargets_Delete_badStatusCode(t *testing.T) {
	stubSetup()
	defer stubTeardown()

	mux.HandleFunc(fmt.Sprintf("/upstreams/%s/targets/i", upstreamName), func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		fmt.Fprint(w, `{"error":"e"}`)
	})

	_, err := client.Targets.Delete(upstreamName, "i")
	if err == nil {
		t.Error("Expected error to be returned")
	}
}

func TestTargets_Post(t *testing.T) {
	stubSetup()
	defer stubTeardown()

	input := sampleTarget()

	mux.HandleFunc(fmt.Sprintf("/upstreams/%s/targets", upstreamName), func(w http.ResponseWriter, r *http.Request) {
		v := new(Target)
		json.NewDecoder(r.Body).Decode(v)
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		testMethod(t, r, "POST")
	})

	_, err := client.Targets.Post(upstreamName, input)
	if err != nil {
		t.Errorf("Targets.Post returned error: %v", err)
	}
}

func TestTargets_Post_badStatusCode(t *testing.T) {
	stubSetup()
	defer stubTeardown()

	mux.HandleFunc(fmt.Sprintf("/upstreams/%s/targets", upstreamName), func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		fmt.Fprint(w, `{"error":"e"}`)
	})

	input := sampleTarget()

	_, err := client.Targets.Post(upstreamName, input)
	if err == nil {
		t.Error("Expected error to be returned")
	}
}

func TestTargets_GetAllActive(t *testing.T) {
	stubSetup()
	defer stubTeardown()

	v := &Targets{Total: 1, Next: "n", Data: []*Target{sampleTarget()}}

	mux.HandleFunc(fmt.Sprintf("/upstreams/%s/targets/active", upstreamName), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		json.NewEncoder(w).Encode(v)
	})

	targets, _, err := client.Targets.GetAllActive(upstreamName)
	if err != nil {
		t.Errorf("Targets.GetAllActive returned error: %v", err)
	}

	if !reflect.DeepEqual(targets, v) {
		t.Errorf("Targets.GetAllActive returned %+v, want %+v", targets, v)
	}
}

func TestTargets_GetAllActive_NoneActive(t *testing.T) {
	stubSetup()
	defer stubTeardown()

	v := &Targets{Total: 0, Data: []*Target{}}

	mux.HandleFunc(fmt.Sprintf("/upstreams/%s/targets/active", upstreamName), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		json.NewEncoder(w).Encode(v)
	})

	targets, _, err := client.Targets.GetAllActive(upstreamName)
	if err != nil {
		t.Errorf("Targets.GetAllActive returned error: %v", err)
	}

	if !reflect.DeepEqual(targets, v) {
		t.Errorf("Targets.GetAllActive returned %+v, want %+v", targets, v)
	}
}

func TestTargets_GetAllActive_badStatusCode(t *testing.T) {
	stubSetup()
	defer stubTeardown()

	mux.HandleFunc(fmt.Sprintf("/upstreams/%s/targets/active", upstreamName), func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		fmt.Fprint(w, `{"error":"e"}`)
	})

	_, _, err := client.Targets.GetAllActive(upstreamName)
	if err == nil {
		t.Error("Expected error to be returned")
	}
}

func sampleTarget() *Target {
	return &Target{
		Target: "service:80",
		Weight: 34,
	}

}
