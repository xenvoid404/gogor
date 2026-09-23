package gogor

import (
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
