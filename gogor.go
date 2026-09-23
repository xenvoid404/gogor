package gogor

import (
	"net/http"
	"strings"
	"time"
)

type Client struct {
	httpClient *http.Client
	baseURL    string
	timeout    time.Duration
}

func New() *Client {
	return &Client{
		httpClient: &http.Client{},
		timeout:    15 * time.Second,
	}
}

func (c *Client) BaseURL(url string) *Client {
	c.baseURL = strings.TrimRight(url, "/")
	return c
}

func (c *Client) Timeout(duration time.Duration) *Client {
	c.timeout = duration
	c.httpClient.Timeout = duration
	return c
}
