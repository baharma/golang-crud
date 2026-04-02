# Golang CRUD

Project ini adalah CRUD service berbasis Go, Gin, GORM, dan MySQL.

Dokumentasi penjelasan alur non-code yang lebih detail tersedia di [ALUR_APLIKASI.md](./ALUR_APLIKASI.md).

Saat ini project punya tiga entrypoint utama:

- `cmd/web`: menjalankan HTTP API
- `cmd/migrate`: membuat database jika belum ada lalu menjalankan migrasi schema
- `cmd/worker`: disiapkan untuk background job, tetapi saat ini masih placeholder

## Struktur Singkat

```text
cmd/
  web/       HTTP server
  migrate/   database migration command
  worker/    background worker placeholder

internal/
  config/    pembacaan konfigurasi
  database/  koneksi database dan migrasi
  entity/    model GORM
```

## Kebutuhan

Sebelum menjalankan project, siapkan:

- Go
- MySQL yang aktif
- file `.env`

Versi Go yang saat ini dipakai di mesin pengembangan adalah `go1.26.1`.

## Konfigurasi `.env`

Isi minimal yang dibutuhkan:

```env
APP_PORT=8080
DB_HOST=127.0.0.1
DB_PORT=3306
DB_NAME=crud_go
DB_USER=root
DB_PASS=your_password
```

Keterangan:

- `APP_PORT`: port untuk HTTP server
- `DB_HOST`: host MySQL
- `DB_PORT`: port MySQL
- `DB_NAME`: nama database aplikasi
- `DB_USER`: username MySQL
- `DB_PASS`: password MySQL

## Penjelasan Command

### `cmd/web`

Command ini menjalankan HTTP server utama.

Tanggung jawabnya:

- membaca konfigurasi dari `.env`
- memastikan database bisa diakses
- menyalakan server Gin
- menyediakan endpoint health check di `/health`

Jalankan dengan:

```bash
go run ./cmd/web
```

### `cmd/migrate`

Command ini dipakai untuk migrasi database.

Tanggung jawabnya:

- konek ke server MySQL
- membuat database `crud_go` jika belum ada
- menjalankan `AutoMigrate` untuk model:
  - `categories`
  - `books`

Jalankan dengan:

```bash
go run ./cmd/migrate
```

Kenapa dipisah dari `web`:

- startup API jadi lebih bersih
- migrasi bisa dijalankan manual saat dibutuhkan
- lebih mudah dipakai di deployment atau CI/CD

### `cmd/worker`

`Worker` adalah proses background terpisah dari web server.

Biasanya worker dipakai untuk:

- job queue
- kirim email atau notifikasi
- sinkronisasi data
- scheduled task
- proses berat yang tidak cocok dijalankan langsung di request HTTP

Di project ini, worker belum memiliki implementasi nyata dan masih placeholder.

## Cara Menjalankan

1. Pastikan MySQL aktif.
2. Pastikan `.env` terisi benar.
3. Jalankan migrasi:

```bash
go run ./cmd/migrate
```

4. Jalankan web server:

```bash
go run ./cmd/web
```

5. Cek health endpoint:

```bash
curl http://127.0.0.1:8080/health
```

Jika berhasil, response-nya akan seperti ini:

```json
{"message":"CRUD API is running","status":"healthy"}
```

## Penjelasan Migrasi yang Dibuat

Migrasi yang baru ditambahkan bekerja seperti ini:

1. membaca konfigurasi database dari `.env`
2. konek ke MySQL tanpa memilih database aplikasi
3. membuat database jika belum ada
4. konek ke database aplikasi
5. menjalankan `AutoMigrate` untuk membuat atau menyesuaikan tabel

Schema yang saat ini dimigrasikan:

- `categories`
- `books`

Catatan:

- field string pada entity sudah diberi ukuran `255` agar kompatibel dengan index MySQL
- `cmd/web` sekarang tidak lagi menjalankan migrasi otomatis saat start

## File Penting

- `cmd/web/main.go`
- `cmd/migrate/main.go`
- `cmd/worker/main.go`
- `internal/database/connection.go`
- `internal/database/migrate.go`
- `internal/entity/category.go`
- `internal/entity/book.go`

## Status Saat Ini

Yang sudah diverifikasi:

- `go test ./...` berhasil
- `go run ./cmd/migrate` berhasil
- `go run ./cmd/web` berhasil
- endpoint `/health` merespons normal

Yang belum ada:

- endpoint CRUD untuk data
- implementasi worker yang nyata
- migration versioning terpisah selain `AutoMigrate`
