package kong

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestConsumersService_Get(t *testing.T) {
	stubSetup()
	defer stubTeardown()

	mux.HandleFunc("/consumers/i", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":"i"}`)
	})

	consumer, _, err := client.Consumers.Get("i")
	if err != nil {
		t.Errorf("Consumers.Get returned error: %v", err)
	}

	want := &Consumer{ID: "i"}
	if !reflect.DeepEqual(consumer, want) {
		t.Errorf("Consumers.Get returned %+v, want %+v", consumer, want)
	}
}

func TestConsumersService_Get_invalidApi(t *testing.T) {
	_, _, err := client.Consumers.Get("%")
	testURLParseError(t, err)
}

func TestConsumersService_Get_badStatusCode(t *testing.T) {
	stubSetup()
	defer stubTeardown()

	mux.HandleFunc("/consumers/i", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		fmt.Fprint(w, `{"error":"e"}`)
	})

	_, _, err := client.Consumers.Get("i")
	if err == nil {
		t.Error("Expected error to be returned")
	}
}

func TestConsumersService_Patch_byUsername(t *testing.T) {
	stubSetup()
	defer stubTeardown()

	input := &Consumer{Username: "u"}

	mux.HandleFunc("/consumers/u", func(w http.ResponseWriter, r *http.Request) {
		v := new(Consumer)
		json.NewDecoder(r.Body).Decode(v)
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		testMethod(t, r, "PATCH")

	})

	_, err := client.Consumers.Patch(input)
	if err != nil {
		t.Errorf("Consumers.Patch returned error: %v", err)
	}
}

func TestConsumersService_Patch_byID(t *testing.T) {
	stubSetup()
	defer stubTeardown()

	input := &Consumer{ID: "i"}

	mux.HandleFunc("/consumers/i", func(w http.ResponseWriter, r *http.Request) {
		v := new(Consumer)
		json.NewDecoder(r.Body).Decode(v)
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		testMethod(t, r, "PATCH")

	})

	_, err := client.Consumers.Patch(input)
	if err != nil {
		t.Errorf("Consumers.Patch returned error: %v", err)
	}
}

func TestConsumersService_Patch_missingUsernameOrID(t *testing.T) {
	input := &Consumer{CustomID: "c"}
	_, err := client.Consumers.Patch(input)
	if err == nil {
		t.Error("Expected error to be returned")
	}
}

func TestConsumersService_Patch_invalidConsumer(t *testing.T) {
	input := &Consumer{ID: "%"}
	_, err := client.Consumers.Patch(input)
	testURLParseError(t, err)
}

func TestConsumersService_Patch_badStatusCode(t *testing.T) {
	stubSetup()
	defer stubTeardown()

	mux.HandleFunc("/consumers/i", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		fmt.Fprint(w, `{"error":"e"}`)
	})

	input := &Consumer{ID: "i"}

	_, err := client.Consumers.Patch(input)
	if err == nil {
		t.Error("Expected error to be returned")
	}
}

func TestConsumersService_Delete(t *testing.T) {
	stubSetup()
	defer stubTeardown()

	mux.HandleFunc("/consumers/i", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	_, err := client.Consumers.Delete("i")
	if err != nil {
		t.Errorf("Consumers.Delete returned error: %v", err)
	}
}

func TestConsumersService_Delete_invalidConsumer(t *testing.T) {
	_, err := client.Consumers.Delete("%")
	testURLParseError(t, err)
}

func TestConsumersService_Delete_badStatusCode(t *testing.T) {
	stubSetup()
	defer stubTeardown()

	mux.HandleFunc("/consumers/i", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		fmt.Fprint(w, `{"error":"e"}`)
	})

	_, err := client.Consumers.Delete("i")
	if err == nil {
		t.Error("Expected error to be returned")
	}
}

func TestConsumersService_Post(t *testing.T) {
	stubSetup()
	defer stubTeardown()

	input := &Consumer{ID: "i"}

	mux.HandleFunc("/consumers", func(w http.ResponseWriter, r *http.Request) {
		v := new(Consumer)
		json.NewDecoder(r.Body).Decode(v)
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		testMethod(t, r, "POST")

	})

	_, err := client.Consumers.Post(input)
	if err != nil {
		t.Errorf("Apis.Post returned error: %v", err)
	}
}

func TestConsumersService_Post_badStatusCode(t *testing.T) {
	stubSetup()
	defer stubTeardown()

	mux.HandleFunc("/consumers", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		fmt.Fprint(w, `{"error":"e"}`)
	})

	input := &Consumer{ID: "i"}

	_, err := client.Consumers.Post(input)
	if err == nil {
		t.Error("Expected error to be returned")
	}
}

func TestConsumersService_GetAll(t *testing.T) {
	stubSetup()
	defer stubTeardown()

	v := &Consumers{Total: 1, Next: "n", Data: []Consumer{{ID: "i"}}}

	mux.HandleFunc("/consumers", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"offset": "o", "custom_id": "c"})
		json.NewEncoder(w).Encode(v)
	})

	opt := &ConsumersGetAllOptions{Offset: "o", CustomID: "c"}
	consumers, _, err := client.Consumers.GetAll(opt)
	if err != nil {
		t.Errorf("Consumers.GetAll returned error: %v", err)
	}

	want := &Consumers{Total: 1, Next: "n", Data: []Consumer{{ID: "i"}}}
	if !reflect.DeepEqual(consumers, want) {
		t.Errorf("Apis.GetAll returned %+v, want %+v", consumers, want)
	}
}

func TestConsumersService_GetAll_badStatusCode(t *testing.T) {
	stubSetup()
	defer stubTeardown()

	mux.HandleFunc("/consumers", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		fmt.Fprint(w, `{"error":"e"}`)
	})

	_, _, err := client.Consumers.GetAll(nil)
	if err == nil {
		t.Error("Expected error to be returned")
	}
}
