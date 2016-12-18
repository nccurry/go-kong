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

func TestPluginsService_Patch_invalidPlugin(t *testing.T) {
	input := &Plugin{ID: "%"}
	_, err := client.Plugins.Patch("a", input)
	testURLParseError(t, err)
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

func TestPluginsService_Delete_invalidPlugin(t *testing.T) {
	_, err := client.Plugins.Delete("a", "%")
	testURLParseError(t, err)
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

func TestPluginsService_GetSchema(t *testing.T) {
	stubSetup()
	defer stubTeardown()

	mux.HandleFunc("/plugins/schema/s", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"f":"f"}`)
	})

	schemas, _, err := client.Plugins.GetSchema("s")
	if err != nil {
		t.Errorf("Plugins.Get returned error: %v", err)
	}

	want := make(map[string]interface{})
	want["f"] = "f"
	if !reflect.DeepEqual(schemas, want) {
		t.Errorf("Plugins.GetEnabled returned %+v, want %+v", schemas, want)
	}
}

func TestPluginsService_GetSchema_invalidSchema(t *testing.T) {
	_, _, err := client.Plugins.GetSchema("%")
	testURLParseError(t, err)
}

func TestPluginsService_GetSchema_badStatusCode(t *testing.T) {
	stubSetup()
	defer stubTeardown()

	mux.HandleFunc("/plugins/schema/s", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		fmt.Fprint(w, `{"error":"e"}`)
	})

	_, _, err := client.Plugins.GetSchema("s")
	if err == nil {
		t.Error("Expected error to be returned")
	}
}

func TestPluginsService_IsZero(t *testing.T) {
	var f func()
	if !isZero(reflect.ValueOf(f)) {
		t.Error("Expected true for 'zero' function")
	}
	var m map[string]interface{}
	if !isZero(reflect.ValueOf(m)) {
		t.Error("Expected true for 'zero' map")
	}
	var s []string
	if !isZero(reflect.ValueOf(s)) {
		t.Error("Expected true for 'zero' slice")
	}
	var a [1]string
	if !isZero(reflect.ValueOf(a)) {
		t.Error("Expected true for 'zero' array")
	}
	var x struct{}
	if !isZero(reflect.ValueOf(x)) {
		t.Error("Expected true for 'zero' struct")
	}
	p := new(struct{})
	if !isZero(reflect.ValueOf(p)) {
		t.Error("Expected true for 'zero' pointer")
	}
}

func TestPluginsService_ToMap(t *testing.T) {
	type T struct {
		F1 string `json:"f_1,omitempty"`
		F2 int    `json:"f_2,omitempty"`
	}
	v := &T{F1: "f1"}
	got := ToMap(v)

	want := make(map[string]interface{})
	want["f_1"] = "f1"
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ToMap returned %+v, want %+v", got, want)
	}
}
