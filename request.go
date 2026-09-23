package gogor

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// DefaultTimeout adalah batas waktu request yang dipakai kalau Timeout tidak diatur.
const DefaultTimeout = 15 * time.Second

// Request adalah satu HTTP request yang disusun secara step-by-step lewat
// method Set* sebelum dieksekusi dengan method terminal seperti Get atau Post.
//
// Request tidak aman dipakai concurrent dari banyak goroutine sekaligus;
// buat instance baru lewat New() untuk setiap request.
type Request struct {
	client  *http.Client
	ctx     context.Context
	url     string
	headers map[string]string
	params  map[string]string
	body    any
	timeout time.Duration
}

// New membuat Request baru dengan context.Background() dan timeout default
// (DefaultTimeout).
func New() *Request {
	return &Request{
		client:  sharedClient,
		ctx:     context.Background(),
		timeout: DefaultTimeout,
	}
}

// Context mengganti context dasar untuk request ini. Kalau ctx nil,
// context.Background() dipakai sebagai gantinya.
func (r *Request) Context(ctx context.Context) *Request {
	if ctx == nil {
		ctx = context.Background()
	}
	r.ctx = ctx
	return r
}

// Timeout mengatur batas waktu request ini, diterapkan di atas Context lewat
// context.WithTimeout. Nilai nol atau negatif berarti tidak ada batas waktu
// tambahan dari sisi Request — request hanya dibatasi oleh Context yang
// sedang dipakai (kalau ada deadline-nya).
func (r *Request) Timeout(d time.Duration) *Request {
	r.timeout = d
	return r
}

// URL mengatur URL tujuan request. Harus berupa URL absolut, diawali
// "http://" atau "https://".
func (r *Request) URL(url string) *Request {
	r.url = url
	return r
}

// Header menambahkan satu header untuk request ini. Pemanggilan berulang
// dengan key yang sama akan menimpa nilai sebelumnya.
func (r *Request) Header(k, v string) *Request {
	if r.headers == nil {
		r.headers = make(map[string]string)
	}
	r.headers[k] = v
	return r
}

// Bearer adalah alias singkat untuk Header("Authorization", "Bearer "+token).
func (r *Request) Bearer(token string) *Request {
	return r.Header("Authorization", "Bearer "+token)
}

// Query menambahkan satu query parameter untuk request ini. Pemanggilan
// berulang dengan key yang sama akan menimpa nilai sebelumnya.
func (r *Request) Query(k, v string) *Request {
	if r.params == nil {
		r.params = make(map[string]string)
	}
	r.params[k] = v
	return r
}

// Body mengatur body request. Cara encoding tergantung tipe v:
//   - nil, io.Reader: dikirim apa adanya, tanpa Content-Type otomatis
//   - []byte: dikirim apa adanya (Content-Type: application/octet-stream)
//   - string: dikirim apa adanya (Content-Type: text/plain)
//   - url.Values: di-encode sebagai form (Content-Type: application/x-www-form-urlencoded)
//   - selain itu: di-encode sebagai JSON (Content-Type: application/json)
func (r *Request) Body(v any) *Request {
	r.body = v
	return r
}

// Get mengeksekusi request ini sebagai HTTP GET.
func (r *Request) Get() (*Response, error) {
	return r.do(http.MethodGet)
}

// Post mengeksekusi request ini sebagai HTTP POST.
func (r *Request) Post() (*Response, error) {
	return r.do(http.MethodPost)
}

// Put mengeksekusi request ini sebagai HTTP PUT.
func (r *Request) Put() (*Response, error) {
	return r.do(http.MethodPut)
}

// Patch mengeksekusi request ini sebagai HTTP PATCH.
func (r *Request) Patch() (*Response, error) {
	return r.do(http.MethodPatch)
}

// Delete mengeksekusi request ini sebagai HTTP DELETE.
func (r *Request) Delete() (*Response, error) {
	return r.do(http.MethodDelete)
}

// Head mengeksekusi request ini sebagai HTTP HEAD.
func (r *Request) Head() (*Response, error) {
	return r.do(http.MethodHead)
}

// do menyusun dan menjalankan HTTP request, lalu mengembalikan hasilnya
// sebagai *Response. do tidak menganggap status code >= 400 sebagai error;
// error hanya dikembalikan untuk kegagalan di level pembuatan request atau
// jaringan (URL tidak valid, gagal konek, gagal baca body, dst).
func (r *Request) do(method string) (*Response, error) {
	if r.url == "" {
		return nil, fmt.Errorf("gogor: URL wajib di-isi")
	}
	if !strings.HasPrefix(r.url, "http://") && !strings.HasPrefix(r.url, "https://") {
		return nil, fmt.Errorf("gogor: URL tidak valid %q, harus diawali http:// atau https://", r.url)
	}

	fullURL, err := r.buildURL()
	if err != nil {
		return nil, err
	}

	bodyReader, contentType, err := r.buildBody()
	if err != nil {
		return nil, err
	}

	// Timeout hanya diterapkan kalau > 0. Nilai nol atau negatif berarti
	// request mengikuti Context apa adanya, tanpa batas waktu tambahan —
	// context.WithTimeout(ctx, 0) akan membuat context langsung expired,
	// jadi kasus ini harus ditangani terpisah.
	ctx := r.ctx
	if r.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, r.timeout)
		defer cancel()
	}

	httpReq, err := http.NewRequestWithContext(ctx, method, fullURL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("gogor: gagal build request: %w", err)
	}

	if contentType != "" {
		httpReq.Header.Set("Content-Type", contentType)
	}
	for k, v := range r.headers {
		httpReq.Header.Set(k, v)
	}

	httpResp, err := r.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("gogor: request gagal: %w", err)
	}
	defer func() {
		_ = httpResp.Body.Close()
	}()

	rawBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("gogor: gagal baca response: %w", err)
	}

	return &Response{
		Status:     httpResp.StatusCode,
		StatusText: httpResp.Status,
		Headers:    httpResp.Header,
		Body:       rawBody,
	}, nil
}

// buildURL menyisipkan query parameter (jika ada) ke dalam URL request.
func (r *Request) buildURL() (string, error) {
	if len(r.params) == 0 {
		return r.url, nil
	}

	u, err := url.Parse(r.url)
	if err != nil {
		return "", fmt.Errorf("gogor: gagal parse url %q: %w", r.url, err)
	}

	q := u.Query()
	for k, v := range r.params {
		q.Set(k, v)
	}

	u.RawQuery = q.Encode()
	return u.String(), nil
}

// buildBody menentukan reader body request dan content-type-nya berdasarkan
// tipe data yang diatur lewat Body. Lihat dokumentasi Body untuk aturan
// encoding per tipe.
func (r *Request) buildBody() (io.Reader, string, error) {
	if r.body == nil {
		return nil, "", nil
	}

	switch v := r.body.(type) {
	case io.Reader:
		return v, "", nil
	case []byte:
		return bytes.NewReader(v), "application/octet-stream", nil
	case string:
		return strings.NewReader(v), "text/plain", nil
	case url.Values:
		return strings.NewReader(v.Encode()), "application/x-www-form-urlencoded", nil
	default:
		b, err := json.Marshal(v)
		if err != nil {
			return nil, "", fmt.Errorf("gogor: gagal marshal body: %w", err)
		}
		return bytes.NewReader(b), "application/json", nil
	}
}
