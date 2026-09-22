package gogor

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGet(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Castella"))
	}))
	defer srv.Close()

	body, err := Get(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	if body != "Castella" {
		t.Fatalf("body harusnya Castella, ini malah %s", body)
	}
}
