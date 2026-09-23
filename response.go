package gogor

import (
	"encoding/json"
	"net/http"
)

// Response adalah hasil dari sebuah HTTP request yang sudah selesai dieksekusi.
type Response struct {
	// Status adalah kode status HTTP (contoh: 200).
	Status int

	// StatusText adalah baris status HTTP (contoh: "200 OK").
	StatusText string

	// Headers adalah header dari response.
	Headers http.Header

	// Body adalah body response dalam bentuk mentah ([]byte).
	Body []byte
}

// JSON men-decode Body ke dalam v. v harus berupa pointer, mengikuti aturan
// encoding/json.Unmarshal.
//
// JSON tidak melakukan pengecekan status code, jadi bisa dipakai baik untuk
// body sukses maupun body error — tinggal sesuaikan struct yang kamu berikan
// sebagai v dengan bentuk response dari masing-masing kasus.
func (r *Response) JSON(v any) error {
	return json.Unmarshal(r.Body, v)
}
