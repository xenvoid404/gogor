package gogor

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGet(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Castella"))
	}))
	defer srv.Close()

	client := New(srv.URL)
	body, err := client.Get("/h2h")
	if err != nil {
		t.Fatal(err)
	}
	if body != "Castella" {
		t.Fatalf("body harusnya Castella, ini malah %s", body)
	}
}
