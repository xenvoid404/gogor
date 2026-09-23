// Package gogor adalah wrapper HTTP client minimalis untuk Go, dibangun murni
// di atas standard library net/http. Didesain dengan gaya step-by-step
// (bukan method-chaining generik), supaya request disusun secara
// eksplisit dan gampang dibaca: New().URL(...).Header(...).Get().
//
// Contoh pemakaian dasar:
//
//	resp, err := gogor.New().
//	    URL("https://api.example.com/users/1").
//	    Get()
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	var user User
//	if err := resp.JSON(&user); err != nil {
//	    log.Fatal(err)
//	}
//
// gogor tidak menganggap status code >= 400 sebagai error Go, sama seperti
// perilaku http.Client.Do bawaan stdlib. Caller bertanggung jawab mengecek
// resp.Status sesuai kebutuhan masing-masing.
package gogor

import "net/http"

// sharedClient adalah http.Client yang dipakai bersama oleh semua Request.
// http.Client aman dipakai secara concurrent, dan connection pooling di
// dalamnya (lewat Transport) memang didesain untuk dipakai berulang kali,
// bukan dibuat baru di setiap request.
var sharedClient = &http.Client{}
