<p align="center">
  <img src="logo.png" alt="gogor">
</p>

<p align="center">
  <strong>gogor</strong> — HTTP client minimalis untuk Go. Hanya memakai standard library,
  tanpa dependensi eksternal. Didesain dengan gaya <strong>fluent interface</strong>
  supaya request disusun secara eksplisit dan mudah dibaca.
</p>

<p align="center">
  <a href="https://github.com/xenvoid404/gogor/actions" target="_blank">
    <img src="https://github.com/xenvoid404/gogor/actions/workflows/workflow.yaml/badge.svg" alt="Test Status">
  </a>
  <a href="https://coveralls.io/github/xenvoid404/gogor?branch=master" target="_blank">
    <img src="https://coveralls.io/repos/github/xenvoid404/gogor/badge.svg?branch=master" alt="Coverage Status">
  </a>
  <br />
  <a href="https://pkg.go.dev/github.com/xenvoid404/gogor" target="_blank">
    <img src="https://pkg.go.dev/badge/github.com/xenvoid404/gogor.svg" alt="Go Reference">
  </a>
  <a href="https://opensource.org/licenses/MIT" target="_blank">
    <img src="https://img.shields.io/badge/License-MIT-yellow.svg" alt="License">
  </a>
  <br />
  <a href="https://github.com/xenvoid404/gogor/releases" target="_blank">
    <img src="https://img.shields.io/github/v/release/xenvoid404/gogor" alt="Release">
  </a>
</p>

---

## Daftar Isi

- [Instalasi](#instalasi)
- [Kenapa gogor](#kenapa-gogor)
- [Quickstart](#quickstart)
- [Filosofi](#filosofi)
- [Dokumentasi](#dokumentasi)
- [Kontribusi](#kontribusi)
- [Lisensi](#lisensi)

## Instalasi

```bash
go get github.com/xenvoid404/gogor
```

## Kenapa gogor

- **Zero dependency** — cuma pakai `net/http`, `encoding/json`, dan paket standard library lainnya.
- **Eksplisit** — tidak ada magic, tidak ada interceptor tersembunyi. Apa yang kamu tulis itu yang jalan.
- **Perilaku sesuai stdlib** — status code `>= 400` **tidak** dianggap error Go, persis seperti `http.Client.Do` bawaan. Error cuma dikembalikan untuk kegagalan level request/jaringan (URL tidak valid, gagal konek, timeout, dll).
- **`resp.Ok`** — pengecekan status sukses satu baris, tanpa harus menghafal rentang kode HTTP tiap kali.

## Quickstart

**GET dengan query parameter:**

```go
resp, err := gogor.New().
    URL("https://api.example.com/users").
    Query("page", "2").
    Query("limit", "15").
    Get()
```

**POST dengan body JSON otomatis:**

```go
resp, err := gogor.New().
    URL("https://api.example.com/users").
    Body(User{Name: "Castella"}).
    Post()
```

**Header, Bearer token, dan timeout custom:**

```go
resp, err := gogor.New().
    URL("https://api.example.com/me").
    Bearer("token123").
    Header("X-Request-Id", "abc").
    Timeout(5 * time.Second).
    Get()
```

**Cek status dan decode response:**

```go
resp, err := gogor.New().URL("https://api.example.com/users/1").Get()
if err != nil {
    log.Fatal(err)
}

if !resp.Ok {
    log.Printf("request gagal: %d %s", resp.Status, resp.StatusText)
    return
}

var result MyStruct
if err := resp.JSON(&result); err != nil {
    log.Fatal(err)
}
```

> `err` hanya diisi untuk kegagalan level request/jaringan (URL tidak valid, gagal konek, timeout). Status code selalu dicek lewat `resp.Ok` atau `resp.Status`, bukan lewat `err`.

## Filosofi

gogor sengaja dibuat kecil. Setiap fitur yang mau ditambahkan harus lolos satu pertanyaan: **apakah ini masih minimalis?**

Kalau sebuah fitur bisa dilakukan pemanggil dengan beberapa baris kode sendiri, fitur itu tidak akan masuk ke gogor. Kalau fitur itu butuh implementasi yang kompleks dan rawan salah — barulah masuk akal untuk dipertimbangkan.

Saat ini gogor **tidak** menyertakan:

- Retry otomatis
- Middleware/interceptor
- Circuit breaker
- Response caching

Semua itu bisa dibangun di atas gogor sesuai kebutuhan aplikasimu masing-masing — misalnya lewat `http.RoundTripper` kustom — atau pakai library yang lebih lengkap seperti [resty](https://github.com/go-resty/resty) kalau memang dibutuhkan.

## Dokumentasi

- [Dokumentasi Wiki](https://github.com/xenvoid404/gogor/wiki) — dokumentasi lengkap
- [Referensi API](https://pkg.go.dev/github.com/xenvoid404/gogor) — referensi API otomatis (dari doc comment di kode)

## Kontribusi

Issue dan pull request dipersilakan. Sebelum mengirim PR:

1. Jalankan `go test ./... -race` dan pastikan semua lulus.
2. Jalankan `go vet ./...`.
3. Ikuti gaya kode dan dokumentasi (bahasa Indonesia semi-formal) yang sudah ada.

## Lisensi

Didistribusikan di bawah lisensi [MIT](LICENSE).
