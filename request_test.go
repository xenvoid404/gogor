package gogor

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func TestRequest_buildURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		params  map[string]string
		want    string
		wantErr bool
	}{
		{
			name: "tanpa params",
			url:  "https://api.example.com/users",
			want: "https://api.example.com/users",
		},
		{
			name:   "dengan params",
			url:    "https://api.example.com/users",
			params: map[string]string{"page": "1"},
			want:   "https://api.example.com/users?page=1",
		},
		{
			name:   "params digabung dengan query yang sudah ada di url",
			url:    "https://api.example.com/users?sort=asc",
			params: map[string]string{"page": "1"},
			want:   "https://api.example.com/users?page=1&sort=asc",
		},
		{
			name:    "url tidak valid",
			url:     "https://example.com/%zz",
			params:  map[string]string{"a": "b"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Request{url: tt.url, params: tt.params}
			got, err := r.buildURL()

			if tt.wantErr {
				if err == nil {
					t.Fatal("harusnya error, tapi malah nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("error tidak diketahui: %v", err)
			}
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestRequest_buildBody(t *testing.T) {
	type jsonPayload struct {
		Name string `json:"name"`
	}

	tests := []struct {
		name            string
		body            any
		wantContentType string
		wantData        string
		wantErr         bool
	}{
		{
			name: "nil",
		},
		{
			name:            "string",
			body:            "castella",
			wantContentType: "text/plain",
			wantData:        "castella",
		},
		{
			name:            "bytes",
			body:            []byte("castella"),
			wantContentType: "application/octet-stream",
			wantData:        "castella",
		},
		{
			name:            "io.Reader dilewatkan apa adanya",
			body:            bytes.NewReader([]byte("castella")),
			wantContentType: "",
			wantData:        "castella",
		},
		{
			name:            "url.Values di-encode sebagai form",
			body:            url.Values{"a": []string{"1"}},
			wantContentType: "application/x-www-form-urlencoded",
			wantData:        "a=1",
		},
		{
			name:            "struct di-encode sebagai json",
			body:            jsonPayload{Name: "castella"},
			wantContentType: "application/json",
			wantData:        `{"name":"castella"}`,
		},
		{
			name:    "tipe yang gagal di-marshal json",
			body:    make(chan int),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Request{body: tt.body}
			reader, contentType, err := r.buildBody()

			if tt.wantErr {
				if err == nil {
					t.Fatal("harusnya error, malah nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("error tidak diketahui: %v", err)
			}
			if contentType != tt.wantContentType {
				t.Errorf("content-type = %q, harusnya %q", contentType, tt.wantContentType)
			}

			if tt.body == nil {
				if reader != nil {
					t.Error("harusnya nil reader untuk body nil")
				}
				return
			}

			got, err := io.ReadAll(reader)
			if err != nil {
				t.Fatalf("gagal baca reader: %v", err)
			}
			if string(got) != tt.wantData {
				t.Errorf("data = %q, harusnya %q", got, tt.wantData)
			}
		})
	}
}

func TestRequest_URLValidation(t *testing.T) {
	t.Run("url kosong", func(t *testing.T) {
		_, err := New().Get()
		if err == nil {
			t.Fatal("harusnya error untuk url kosong")
		}
	})

	t.Run("url tanpa scheme", func(t *testing.T) {
		_, err := New().URL("api.example.com/users").Get()
		if err == nil {
			t.Fatal("harusnya error untuk url tanpa http:// atau https://")
		}
	})
}

func TestRequest_Do_EndToEnd(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/echo":
			_, _ = io.Copy(w, r.Body)
		case "/json":
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]string{"hello": "world"})
		case "/headers":
			w.Header().Set("X-Got-Auth", r.Header.Get("Authorization"))
		case "/query":
			_, _ = w.Write([]byte(r.URL.RawQuery))
		case "/status":
			w.WriteHeader(http.StatusTeapot)
		}
	}))
	defer srv.Close()

	t.Run("GET dan decode JSON response", func(t *testing.T) {
		resp, err := New().URL(srv.URL + "/json").Get()
		if err != nil {
			t.Fatalf("error asing: %v", err)
		}
		if resp.Status != http.StatusOK {
			t.Fatalf("status = %d, harusnya %d", resp.Status, http.StatusOK)
		}

		var body map[string]string
		if err := resp.JSON(&body); err != nil {
			t.Fatalf("gagal decode json: %v", err)
		}
		if body["hello"] != "world" {
			t.Errorf("body[hello] = %q, harusnya %q", body["hello"], "world")
		}
	})

	t.Run("POST mengirim body string apa adanya", func(t *testing.T) {
		resp, err := New().URL(srv.URL + "/echo").Body("hello world").Post()
		if err != nil {
			t.Fatalf("error asing: %v", err)
		}
		if string(resp.Body) != "hello world" {
			t.Errorf("body = %q, harusnya %q", resp.Body, "hello world")
		}
	})

	t.Run("Bearer menambahkan header Authorization", func(t *testing.T) {
		resp, err := New().URL(srv.URL + "/headers").Bearer("token123").Get()
		if err != nil {
			t.Fatalf("error asing: %v", err)
		}
		if got := resp.Headers.Get("X-Got-Auth"); got != "Bearer token123" {
			t.Errorf("Authorization = %q, harusnya %q", got, "Bearer token123")
		}
	})

	t.Run("Query menambahkan query parameter", func(t *testing.T) {
		resp, err := New().URL(srv.URL+"/query").Query("page", "2").Get()
		if err != nil {
			t.Fatalf("error asing: %v", err)
		}
		if string(resp.Body) != "page=2" {
			t.Errorf("query = %q, harusnya %q", resp.Body, "page=2")
		}
	})

	t.Run("status >= 400 tidak dianggap error", func(t *testing.T) {
		resp, err := New().URL(srv.URL + "/status").Get()
		if err != nil {
			t.Fatalf("harusnya ngga error, malah error %v", err)
		}
		if resp.Status != http.StatusTeapot {
			t.Errorf("status = %d, harusnya %d", resp.Status, http.StatusTeapot)
		}
	})
}

func TestRequest_Timeout(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(50 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	t.Run("timeout terlampaui mengembalikan error", func(t *testing.T) {
		_, err := New().URL(srv.URL).Timeout(10 * time.Millisecond).Get()
		if err == nil {
			t.Fatal("harusnya timeout error, tapi malah nil")
		}
	})

	t.Run("timeout nol tidak langsung gagal", func(t *testing.T) {
		resp, err := New().URL(srv.URL).Timeout(0).Get()
		if err != nil {
			t.Fatalf("error asing: %v", err)
		}
		if resp.Status != http.StatusOK {
			t.Errorf("status = %d, harusnya %d", resp.Status, http.StatusOK)
		}
	})
}
