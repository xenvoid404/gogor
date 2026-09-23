package gogor

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	client := New()
	client.BaseURL("https://api.example.com")
	client.Timeout(30 * time.Second)

	if client == nil {
		t.Fatal("client nil")
	}
}

func TestGet(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("castella"))
	}))
	defer srv.Close()

	client := New().BaseURL(srv.URL)
	resp, err := client.Request().Get("/ping")
	if err != nil {
		t.Fatal(err)
	}
	if resp.Status != http.StatusOK {
		t.Fatalf("status harusnya 200, ini malah %d", resp.Status)
	}
}
