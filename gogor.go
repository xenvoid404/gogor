package gogor

import (
	"context"
	"io"
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

func (c *Client) Request() *Request {
	return &Request{
		client:  c,
		ctx:     context.Background(),
		baseURL: c.baseURL,
	}
}

type Request struct {
	client  *Client
	ctx     context.Context
	baseURL string
	method  string
	path    string
}

func (r *Request) Context(ctx context.Context) *Request {
	if ctx == nil {
		r.ctx = context.Background()
	}
	r.ctx = ctx
	return r
}

func (r *Request) Get(path string) (*Response, error) {
	r.method = http.MethodGet
	r.path = path
	return r.do()
}

func (r *Request) do() (*Response, error) {
	url := r.baseURL + r.path

	httpReq, err := http.NewRequestWithContext(r.ctx, r.method, url, nil)
	if err != nil {
		return nil, err
	}

	httpResp, err := r.client.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}

	return &Response{
		Status:  httpResp.StatusCode,
		Headers: httpResp.Header,
		Body:    body,
	}, nil
}

type Response struct {
	Status  int
	Headers http.Header
	Body    []byte
}
