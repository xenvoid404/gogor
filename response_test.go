package gogor

import "testing"

func TestResponse_JSON(t *testing.T) {
	r := &Response{Body: []byte(`{"name":"castella"}`)}

	var v struct {
		Name string `json:"name"`
	}
	if err := r.JSON(&v); err != nil {
		t.Fatalf("error tidak diketahui: %v", err)
	}
	if v.Name != "castella" {
		t.Errorf("Name = %q, harusnya %q", v.Name, "john")
	}
}

func TestResponse_JSON_InvalidBody(t *testing.T) {
	r := &Response{Body: []byte("bukan json")}

	var v map[string]string
	if err := r.JSON(&v); err == nil {
		t.Fatal("harusnya error untuk body yang bukan json valid")
	}
}
