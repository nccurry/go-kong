package kong

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestPluginsService_Get(t *testing.T) {
	stubSetup()
	defer stubTeardown()

	mux.HandleFunc("/plugins/i", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":"i"}`)
	})

	plugin, _, err := client.Plugins.Get("i")
	if err != nil {
		t.Errorf("Plugins.Get returned error: %v", err)
	}

	want := &Plugin{ID: "i"}
	if !reflect.DeepEqual(plugin, want) {
		t.Errorf("Plugins.Get returned %+v, want %+v", plugin, want)
	}
}

func TestPluginsService_Get_invalidPlugin(t *testing.T) {
	_, _, err := client.Plugins.Get("%")
	testURLParseError(t, err)
}

func TestPluginsService_Get_badStatusCode(t *testing.T) {
	stubSetup()
	defer stubTeardown()

	mux.HandleFunc("/plugins/i", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		fmt.Fprint(w, `{"error":"e"}`)
	})

	_, _, err := client.Plugins.Get("i")
	if err == nil {
		t.Error("Expected error to be returned")
	}
}

func TestPluginsService_GetEnabled(t *testing.T) {
	stubSetup()
	defer stubTeardown()

	mux.HandleFunc("/plugins/enabled", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"enabled_plugins":["p"]}`)
	})

	plugins, _, err := client.Plugins.GetEnabled()
	if err != nil {
		t.Errorf("Plugins.Get returned error: %v", err)
	}

	want := &EnabledPlugins{[]string{"p"}}
	if !reflect.DeepEqual(plugins, want) {
		t.Errorf("Plugins.GetEnabled returned %+v, want %+v", plugins, want)
	}
}

func TestPluginsService_GetEnabled_badStatusCode(t *testing.T) {
	stubSetup()
	defer stubTeardown()

	mux.HandleFunc("/plugins/enabled", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		fmt.Fprint(w, `{"error":"e"}`)
	})

	_, _, err := client.Plugins.GetEnabled()
	if err == nil {
		t.Error("Expected error to be returned")
	}
}

func TestPluginsService_Patch(t *testing.T) {
	stubSetup()
	defer stubTeardown()

	input := &Plugin{ID: "i"}

	mux.HandleFunc("/apis/a/plugins/i", func(w http.ResponseWriter, r *http.Request) {
		v := new(Plugin)
		json.NewDecoder(r.Body).Decode(v)
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		testMethod(t, r, "PATCH")

	})

	_, err := client.Plugins.Patch("a", input)
	if err != nil {
		t.Errorf("Plugins.Patch returned error: %v", err)
	}
}

func TestPluginsService_Patch_badStatusCode(t *testing.T) {
	stubSetup()
	defer stubTeardown()

	mux.HandleFunc("/apis/a/plugins/i", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		fmt.Fprint(w, `{"error":"e"}`)
	})

	input := &Plugin{ID: "i"}

	_, err := client.Plugins.Patch("a", input)
	if err == nil {
		t.Error("Expected error to be returned")
	}
}

func TestPluginsService_Delete(t *testing.T) {
	stubSetup()
	defer stubTeardown()

	mux.HandleFunc("/apis/a/plugins/i", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	_, err := client.Plugins.Delete("a", "i")
	if err != nil {
		t.Errorf("Plugins.Delete returned error: %v", err)
	}
}

func TestPluginsService_Delete_badStatusCode(t *testing.T) {
	stubSetup()
	defer stubTeardown()

	mux.HandleFunc("/apis/a/plugins/i", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		fmt.Fprint(w, `{"error":"e"}`)
	})

	_, err := client.Plugins.Delete("a", "i")
	if err == nil {
		t.Error("Expected error to be returned")
	}
}

func TestPluginsService_Post(t *testing.T) {
	stubSetup()
	defer stubTeardown()

	input := &Plugin{ID: "i"}

	mux.HandleFunc("/plugins", func(w http.ResponseWriter, r *http.Request) {
		v := new(Plugin)
		json.NewDecoder(r.Body).Decode(v)
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		testMethod(t, r, "POST")

	})

	_, err := client.Plugins.Post(input)
	if err != nil {
		t.Errorf("Apis.Post returned error: %v", err)
	}
}

func TestPluginsService_Post_badStatusCode(t *testing.T) {
	stubSetup()
	defer stubTeardown()

	mux.HandleFunc("/plugins", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		fmt.Fprint(w, `{"error":"e"}`)
	})

	input := &Plugin{ID: "i"}

	_, err := client.Plugins.Post(input)
	if err == nil {
		t.Error("Expected error to be returned")
	}
}

func TestPluginsService_GetAll(t *testing.T) {
	stubSetup()
	defer stubTeardown()

	v := &Plugins{Total: 1, Next: "n", Data: []Plugin{{ID: "i"}}}

	mux.HandleFunc("/plugins", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"offset": "o", "name": "n"})
		json.NewEncoder(w).Encode(v)
	})

	opt := &PluginsGetAllOptions{Offset: "o", Name: "n"}
	plugins, _, err := client.Plugins.GetAll(opt)
	if err != nil {
		t.Errorf("Plugins.GetAll returned error: %v", err)
	}

	want := &Plugins{Total: 1, Next: "n", Data: []Plugin{{ID: "i"}}}
	if !reflect.DeepEqual(plugins, want) {
		t.Errorf("Plugins.GetAll returned %+v, want %+v", plugins, want)
	}
}

func TestPluginsService_GetAll_badStatusCode(t *testing.T) {
	stubSetup()
	defer stubTeardown()

	mux.HandleFunc("/plugins", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		fmt.Fprint(w, `{"error":"e"}`)
	})

	_, _, err := client.Plugins.GetAll(nil)
	if err == nil {
		t.Error("Expected error to be returned")
	}
}