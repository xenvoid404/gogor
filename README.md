# gogor

Wrapper HTTP client minimalis untuk Go, dibangun murni di atas standard library `net/http`. Didesain dengan gaya *step-by-step* (bukan chaining generik) supaya request disusun secara eksplisit dan gampang dibaca.

```go
resp, err := gogor.New().
    URL("https://api.example.com/users/1").
    Get()
if err != nil {
    log.Fatal(err)
}

var user User
if err := resp.JSON(&user); err != nil {
    log.Fatal(err)
}
```

## Instalasi

```bash
go get github.com/xenvoid404/gogor
```

## Kenapa gogor

- **Zero dependency** — cuma pakai `net/http`, `encoding/json`, dan paket standard library lainnya.
- **Eksplisit** — gak ada magic, gak ada interceptor/middleware tersembunyi. Apa yang kamu tulis itu yang jalan.
- **Perilaku sesuai stdlib** — status code `>= 400` **tidak** dianggap error Go, persis seperti `http.Client.Do` bawaan. Error cuma dikembalikan untuk kegagalan level request/jaringan (URL tidak valid, gagal konek, timeout, dll).

## Quickstart

**GET dengan query parameter:**

```go
resp, err := gogor.New().
    URL("https://api.example.com/users").
    Query("page", "2").
    Get()
```

**POST dengan body JSON otomatis:**

```go
resp, err := gogor.New().
    URL("https://api.example.com/users").
    Body(User{Name: "John"}).
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
if resp.Status >= 400 {
    log.Printf("request gagal: %s", resp.StatusText)
    return
}

var result MyStruct
if err := resp.JSON(&result); err != nil {
    log.Fatal(err)
}
```

## Dokumentasi lengkap

Referensi API lengkap, aturan encoding body per tipe, dan catatan desain ada di [Wiki](https://github.com/xenvoid404/gogor/wiki).

## Lisensi

MIT
