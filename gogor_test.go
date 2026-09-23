package gogor

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGet(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Castella"))
	}))
	defer srv.Close()

	client := New(srv.URL)
	resp, err := client.Get("/h2h")
	if err != nil {
		t.Fatal(err)
	}
	if resp.Status != http.StatusOK {
		t.Fatalf("Status harusnya 200, ini malah %d", resp.Status)
	}
	if string(resp.Data) != "Castella" {
		t.Fatalf("body harusnya Castella, ini malah %s", string(resp.Data))
	}
	if got := resp.Headers.Get("Content-Type"); got != "text/plain" {
		t.Fatalf("header Content-Type harusnya text/plain, ini malah %s", got)
	}
}

func TestPost(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Fatalf("header Content-Type harusnya application/json, ini malah %s", ct)
		}

		var in map[string]string
		_ = json.NewDecoder(r.Body).Decode(&in)
		if in["name"] != "Castella" {
			t.Fatalf("name harusnya Castella, ini malah %v", in["name"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id":1,"name":"Castella"}`))
	}))
	defer srv.Close()

	client := New(srv.URL)
	resp, err := client.Post("/user", map[string]string{"name": "Castella"})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Status != http.StatusOK {
		t.Fatalf("Status harusnya 200, ini malah %d", resp.Status)
	}

	var out struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	if err := resp.JSON(&out); err != nil {
		t.Fatal(err)
	}
	if out.ID != 1 || out.Name != "Castella" {
		t.Fatalf("harusnya %+v", out)
	}
}
